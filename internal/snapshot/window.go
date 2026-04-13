package snapshot

import (
	"errors"
	"sort"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// WindowOptions configures the sliding window extraction.
type WindowOptions struct {
	// Start is the inclusive begin of the time window.
	Start time.Time
	// End is the inclusive end of the time window.
	End time.Time
	// MaxSnapshots caps the number of snapshots returned (0 = unlimited).
	MaxSnapshots int
}

// DefaultWindowOptions returns sensible defaults: last 1 hour, no cap.
func DefaultWindowOptions() WindowOptions {
	now := time.Now().UTC()
	return WindowOptions{
		Start:        now.Add(-1 * time.Hour),
		End:          now,
		MaxSnapshots: 0,
	}
}

// Validate returns an error if the options are logically invalid.
func (o WindowOptions) Validate() error {
	if !o.End.After(o.Start) {
		return errors.New("window: End must be after Start")
	}
	if o.MaxSnapshots < 0 {
		return errors.New("window: MaxSnapshots must be >= 0")
	}
	return nil
}

// WindowResult holds the snapshots that fall within the requested window.
type WindowResult struct {
	Snapshots []WindowSnap
	Start     time.Time
	End       time.Time
	Total     int
}

// WindowSnap pairs a timestamp with its port list.
type WindowSnap struct {
	Timestamp time.Time
	Ports     []scanner.Port
}

// Window extracts snapshots from snaps whose timestamps fall within [opts.Start, opts.End].
// Results are ordered by timestamp ascending and capped by MaxSnapshots.
func Window(snaps []Snapshot, opts WindowOptions) (WindowResult, error) {
	if err := opts.Validate(); err != nil {
		return WindowResult{}, err
	}

	var matched []WindowSnap
	for _, s := range snaps {
		ts := s.Timestamp
		if (ts.Equal(opts.Start) || ts.After(opts.Start)) &&
			(ts.Equal(opts.End) || ts.Before(opts.End)) {
			matched = append(matched, WindowSnap{
				Timestamp: ts,
				Ports:     s.Ports,
			})
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Timestamp.Before(matched[j].Timestamp)
	})

	total := len(matched)
	if opts.MaxSnapshots > 0 && len(matched) > opts.MaxSnapshots {
		matched = matched[:opts.MaxSnapshots]
	}

	return WindowResult{
		Snapshots: matched,
		Start:     opts.Start,
		End:       opts.End,
		Total:     total,
	}, nil
}
