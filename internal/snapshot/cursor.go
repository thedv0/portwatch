package snapshot

import (
	"errors"
	"time"
)

// CursorOptions controls pagination behaviour when iterating snapshots.
type CursorOptions struct {
	PageSize int
	After    time.Time
}

// DefaultCursorOptions returns sensible defaults.
func DefaultCursorOptions() CursorOptions {
	return CursorOptions{
		PageSize: 20,
	}
}

// CursorResult holds a page of snapshots and the next cursor position.
type CursorResult struct {
	Page     []Snapshot
	NextAfter time.Time
	HasMore  bool
}

// Cursor pages through a slice of snapshots ordered by timestamp.
func Cursor(snaps []Snapshot, opts CursorOptions) (CursorResult, error) {
	if opts.PageSize <= 0 {
		return CursorResult{}, errors.New("cursor: PageSize must be positive")
	}

	// filter snapshots after the cursor position
	var filtered []Snapshot
	for _, s := range snaps {
		if s.Timestamp.After(opts.After) {
			filtered = append(filtered, s)
		}
	}

	// sort ascending by timestamp
	for i := 1; i < len(filtered); i++ {
		for j := i; j > 0 && filtered[j].Timestamp.Before(filtered[j-1].Timestamp); j-- {
			filtered[j], filtered[j-1] = filtered[j-1], filtered[j]
		}
	}

	hasMore := len(filtered) > opts.PageSize
	page := filtered
	if hasMore {
		page = filtered[:opts.PageSize]
	}

	var nextAfter time.Time
	if len(page) > 0 {
		nextAfter = page[len(page)-1].Timestamp
	}

	return CursorResult{
		Page:      page,
		NextAfter: nextAfter,
		HasMore:   hasMore,
	}, nil
}
