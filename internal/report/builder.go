package report

import (
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Builder constructs a Report from diff and current state.
type Builder struct {
	now func() time.Time
}

// NewBuilder returns a Builder. Inject a custom clock via WithClock.
func NewBuilder() *Builder {
	return &Builder{now: time.Now}
}

// WithClock replaces the internal clock (useful for testing).
func (b *Builder) WithClock(fn func() time.Time) *Builder {
	b.now = fn
	return b
}

// Build creates a Report given a diff result and total open port count.
func (b *Builder) Build(diff snapshot.DiffResult, totalOpen int) Report {
	return Report{
		Timestamp: b.now(),
		Added:     diff.Added,
		Removed:   diff.Removed,
		Total:     totalOpen,
	}
}
