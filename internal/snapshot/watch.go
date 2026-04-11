package snapshot

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// WatchOptions configures the port watcher behavior.
type WatchOptions struct {
	Interval  time.Duration
	Ports     []scanner.PortState
	OnChange  func(diff Diff)
	OnError   func(err error)
}

// Watcher continuously polls for port changes and emits diffs.
type Watcher struct {
	opts    WatchOptions
	scanner *scanner.Scanner
}

// NewWatcher creates a Watcher with the given options.
func NewWatcher(s *scanner.Scanner, opts WatchOptions) *Watcher {
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}
	if opts.OnChange == nil {
		opts.OnChange = func(Diff) {}
	}
	if opts.OnError == nil {
		opts.OnError = func(error) {}
	}
	return &Watcher{opts: opts, scanner: s}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.opts.Interval)
	defer ticker.Stop()

	prev, err := w.scanner.Scan()
	if err != nil {
		w.opts.OnError(err)
		prev = nil
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current, err := w.scanner.Scan()
			if err != nil {
				w.opts.OnError(err)
				continue
			}
			d := Diff(prev, current)
			if len(d.Added) > 0 || len(d.Removed) > 0 {
				w.opts.OnChange(d)
			}
			prev = current
		}
	}
}
