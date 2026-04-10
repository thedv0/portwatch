package snapshot

import (
	"fmt"
	"time"
)

// RetentionPolicy defines how long snapshots and history entries are kept.
type RetentionPolicy struct {
	// MaxAge is the maximum age of a snapshot file before it is eligible for removal.
	MaxAge time.Duration
	// MaxCount is the maximum number of snapshot files to retain.
	MaxCount int
	// MaxHistoryEntries is the maximum number of history entries to keep.
	MaxHistoryEntries int
}

// DefaultRetentionPolicy returns a RetentionPolicy with sensible defaults.
func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{
		MaxAge:            7 * 24 * time.Hour,
		MaxCount:          100,
		MaxHistoryEntries: 500,
	}
}

// Validate checks that the retention policy fields are valid.
func (r RetentionPolicy) Validate() error {
	if r.MaxAge < 0 {
		return fmt.Errorf("retention max_age must be non-negative")
	}
	if r.MaxCount < 0 {
		return fmt.Errorf("retention max_count must be non-negative")
	}
	if r.MaxHistoryEntries < 0 {
		return fmt.Errorf("retention max_history_entries must be non-negative")
	}
	return nil
}

// Apply runs the cleaner using the policy against the given snapshot directory.
func (r RetentionPolicy) Apply(dir string) error {
	cleaner := NewCleaner(dir, r.MaxAge, r.MaxCount)
	return cleaner.Clean()
}
