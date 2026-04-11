package audit

import (
	"fmt"
	"io"
	"os"
)

// RotatingFileLogger wraps a FileLogger and rotates the underlying file
// according to a RotationPolicy before each write when needed.
type RotatingFileLogger struct {
	path   string
	policy RotationPolicy
	inner  *Logger
}

// NewRotatingFileLogger creates a Logger that automatically rotates the log
// file at path according to policy.
func NewRotatingFileLogger(path string, policy RotationPolicy) (*RotatingFileLogger, error) {
	if err := policy.Validate(); err != nil {
		return nil, fmt.Errorf("rotation config: %w", err)
	}
	w, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	return &RotatingFileLogger{
		path:   path,
		policy: policy,
		inner:  NewLogger(w),
	}, nil
}

// Log checks for rotation, then delegates to the inner Logger.
func (r *RotatingFileLogger) Log(event map[string]any) error {
	if err := r.maybeRotate(); err != nil {
		return err
	}
	return r.inner.Log(event)
}

// Warn logs an event at warn level.
func (r *RotatingFileLogger) Warn(event map[string]any) error {
	if err := r.maybeRotate(); err != nil {
		return err
	}
	return r.inner.Warn(event)
}

// Alert logs an event at alert level.
func (r *RotatingFileLogger) Alert(event map[string]any) error {
	if err := r.maybeRotate(); err != nil {
		return err
	}
	return r.inner.Alert(event)
}

func (r *RotatingFileLogger) maybeRotate() error {
	needs, err := NeedsRotation(r.path, r.policy)
	if err != nil || !needs {
		return err
	}
	if err := Rotate(r.path, r.policy); err != nil {
		return fmt.Errorf("rotating logger: %w", err)
	}
	// Re-open the log file after rotation.
	f, err := openOrCreate(r.path)
	if err != nil {
		return fmt.Errorf("rotating logger: reopen: %w", err)
	}
	r.inner = NewLogger(f)
	return nil
}

func openOrCreate(path string) (io.Writer, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, err
	}
	return f, nil
}
