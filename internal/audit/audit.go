// Package audit provides structured audit logging for portwatch daemon events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an audit event.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Event     string            `json:"event"`
	Details   map[string]string `json:"details,omitempty"`
}

// Logger writes audit entries to an output sink.
type Logger struct {
	w     io.Writer
	clock func() time.Time
}

// NewLogger creates a Logger writing to w. If w is nil, os.Stderr is used.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w, clock: time.Now}
}

// Log writes an audit entry at the given level.
func (l *Logger) Log(level Level, event string, details map[string]string) error {
	e := Entry{
		Timestamp: l.clock(),
		Level:     level,
		Event:     event,
		Details:   details,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

// Info logs an informational audit event.
func (l *Logger) Info(event string, details map[string]string) error {
	return l.Log(LevelInfo, event, details)
}

// Warn logs a warning audit event.
func (l *Logger) Warn(event string, details map[string]string) error {
	return l.Log(LevelWarn, event, details)
}

// Alert logs a critical alert audit event.
func (l *Logger) Alert(event string, details map[string]string) error {
	return l.Log(LevelAlert, event, details)
}
