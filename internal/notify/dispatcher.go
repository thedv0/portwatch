package notify

import (
	"fmt"
	"log"

	"github.com/user/portwatch/internal/alert"
)

// Dispatcher fans an alert event out to one or more channels.
type Dispatcher struct {
	channels []Channel
	logger   *log.Logger
}

// NewDispatcher creates a Dispatcher with the provided channels.
func NewDispatcher(logger *log.Logger, channels ...Channel) *Dispatcher {
	if logger == nil {
		logger = log.New(log.Writer(), "[notify] ", log.LstdFlags)
	}
	return &Dispatcher{
		channels: channels,
		logger:   logger,
	}
}

// Add appends a channel to the dispatcher at runtime.
func (d *Dispatcher) Add(c Channel) {
	d.channels = append(d.channels, c)
}

// Dispatch sends the event to all registered channels.
// It collects all errors and returns a combined error if any channel fails.
func (d *Dispatcher) Dispatch(e alert.Event) error {
	var errs []string
	for _, ch := range d.channels {
		if err := ch.Send(e); err != nil {
			d.logger.Printf("channel %q error: %v", ch.Name(), err)
			errs = append(errs, fmt.Sprintf("%s: %v", ch.Name(), err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dispatch errors: %s", joinStrings(errs))
	}
	return nil
}

func joinStrings(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += "; "
		}
		out += s
	}
	return out
}
