package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes an unexpected port listener detected by the scanner.
type Event struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Protocol  string
	Process   string
	Message   string
}

// Notifier sends alert events to one or more destinations.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier that writes to w.
// If w is nil, os.Stderr is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stderr
	}
	return &Notifier{out: w}
}

// Send formats and writes an Event to the configured writer.
func (n *Notifier) Send(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	_, err := fmt.Fprintf(
		n.out,
		"%s [%s] port=%d proto=%s process=%q msg=%s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port,
		e.Protocol,
		e.Process,
		e.Message,
	)
	return err
}

// NewEvent is a convenience constructor for an alert Event.
func NewEvent(level Level, port int, proto, process, msg string) Event {
	return Event{
		Timestamp: time.Now(),
		Level:     level,
		Port:      port,
		Protocol:  proto,
		Process:   process,
		Message:   msg,
	}
}
