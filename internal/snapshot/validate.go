package snapshot

import (
	"errors"
	"fmt"
)

// ValidationLevel indicates the severity of a validation issue.
type ValidationLevel int

const (
	LevelInfo    ValidationLevel = iota
	LevelWarning ValidationLevel = iota
	LevelError   ValidationLevel = iota
)

// ValidationIssue represents a single validation finding for a port entry.
type ValidationIssue struct {
	Port    int
	Process string
	Level   ValidationLevel
	Message string
}

// ValidationResult holds all issues found during validation.
type ValidationResult struct {
	Issues []ValidationIssue
}

// HasErrors returns true if any error-level issues exist.
func (r *ValidationResult) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Level == LevelError {
			return true
		}
	}
	return false
}

// DefaultValidateOptions returns sensible defaults.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		AllowPIDZero:   false,
		MaxPort:        65535,
		RequireProcess: false,
	}
}

// ValidateOptions controls validation behaviour.
type ValidateOptions struct {
	AllowPIDZero   bool
	MaxPort        int
	RequireProcess bool
}

// Validate checks a slice of ports against the given options and returns a result.
func Validate(ports []PortEntry, opts ValidateOptions) (*ValidationResult, error) {
	if opts.MaxPort <= 0 {
		return nil, errors.New("validate: MaxPort must be positive")
	}
	result := &ValidationResult{}
	for _, p := range ports {
		if p.Port < 0 || p.Port > opts.MaxPort {
			result.Issues = append(result.Issues, ValidationIssue{
				Port:    p.Port,
				Process: p.Process,
				Level:   LevelError,
				Message: fmt.Sprintf("port %d out of valid range [0, %d]", p.Port, opts.MaxPort),
			})
		}
		if !opts.AllowPIDZero && p.PID == 0 {
			result.Issues = append(result.Issues, ValidationIssue{
				Port:    p.Port,
				Process: p.Process,
				Level:   LevelWarning,
				Message: fmt.Sprintf("port %d has PID 0 (unknown process)", p.Port),
			})
		}
		if opts.RequireProcess && p.Process == "" {
			result.Issues = append(result.Issues, ValidationIssue{
				Port:    p.Port,
				Process: p.Process,
				Level:   LevelWarning,
				Message: fmt.Sprintf("port %d has no associated process name", p.Port),
			})
		}
	}
	return result, nil
}
