package snapshot

import (
	"errors"
	"sort"
	"time"
)

// RetainOptions controls which snapshots are kept.
type RetainOptions struct {
	// MaxAge removes snapshots older than this duration. Zero disables.
	MaxAge time.Duration
	// MinCount ensures at least this many snapshots are kept even if older than MaxAge.
	MinCount int
	// MaxCount keeps only the newest N snapshots. Zero disables.
	MaxCount int
}

// DefaultRetainOptions returns sensible defaults.
func DefaultRetainOptions() RetainOptions {
	return RetainOptions{
		MaxAge:   72 * time.Hour,
		MinCount: 1,
		MaxCount: 0,
	}
}

// Validate checks that the options are self-consistent.
func (o RetainOptions) Validate() error {
	if o.MaxAge < 0 {
		return errors.New("retain: MaxAge must not be negative")
	}
	if o.MinCount < 0 {
		return errors.New("retain: MinCount must not be negative")
	}
	if o.MaxCount < 0 {
		return errors.New("retain: MaxCount must not be negative")
	}
	return nil
}

// Retain filters a slice of snapshots according to the given options.
// Snapshots are assumed to be in any order; the result is sorted newest-first.
func Retain(snaps []Snapshot, opts RetainOptions) ([]Snapshot, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	// Sort newest first.
	sorted := make([]Snapshot, len(snaps))
	copy(sorted, snaps)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.After(sorted[j].Timestamp)
	})

	now := time.Now()

	// Apply MaxAge, but always keep at least MinCount.
	var result []Snapshot
	for i, s := range sorted {
		keepByMin := opts.MinCount > 0 && i < opts.MinCount
		tooOld := opts.MaxAge > 0 && now.Sub(s.Timestamp) > opts.MaxAge
		if keepByMin || !tooOld {
			result = append(result, s)
		}
	}

	// Apply MaxCount.
	if opts.MaxCount > 0 && len(result) > opts.MaxCount {
		result = result[:opts.MaxCount]
	}

	return result, nil
}
