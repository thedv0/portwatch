package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
)

// Daemon polls open ports on a fixed interval and fires alerts on violations.
type Daemon struct {
	cfg      *rules.Config
	scanner  *scanner.Scanner
	matcher  *alert.Matcher
	notifier *alert.Notifier
	prev     map[string]scanner.PortState
}

// New creates a Daemon from the supplied configuration.
func New(cfg *rules.Config, n *alert.Notifier) *Daemon {
	return &Daemon{
		cfg:      cfg,
		scanner:  scanner.NewScanner(),
		matcher:  alert.NewMatcher(cfg),
		notifier: n,
		prev:     make(map[string]scanner.PortState),
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	interval, err := time.ParseDuration(d.cfg.Interval)
	if err != nil {
		return err
	}

	log.Printf("portwatch: starting — interval %s", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run an immediate first scan before waiting for the first tick.
	if err := d.tick(); err != nil {
		log.Printf("portwatch: scan error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch: shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		}
	}
}

func (d *Daemon) tick() error {
	states, err := d.scanner.Scan()
	if err != nil {
		return err
	}

	events := d.matcher.Evaluate(states, d.prev)
	for _, ev := range events {
		if err := d.notifier.Send(ev); err != nil {
			log.Printf("portwatch: alert send error: %v", err)
		}
	}

	d.prev = states
	return nil
}
