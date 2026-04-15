package snapshot

import (
	"errors"
	"time"
)

// ShiftOptions configures how timestamps in snapshots are shifted.
type ShiftOptions struct {
	// Offset is the duration to add to each snapshot's timestamp.
	Offset time.Duration
	// ClampToNow prevents shifted timestamps from exceeding the current time.
	ClampToNow bool
	// Clock allows injecting a custom time source for testing.
	Clock func() time.Time
}

// DefaultShiftOptions returns a ShiftOptions with sensible defaults.
func DefaultShiftOptions() ShiftOptions {
	return ShiftOptions{
		Offset:     0,
		ClampToNow: false,
		Clock:      time.Now,
	}
}

// Shift applies a time offset to every snapshot in the provided slice.
// Snapshots are returned in the same order with adjusted timestamps.
// An error is returned if opts.Clock is nil.
func Shift(snaps []Snapshot, opts ShiftOptions) ([]Snapshot, error) {
	if opts.Clock == nil {
		return nil, errors.New("shift: Clock must not be nil")
	}

	now := opts.Clock()
	out := make([]Snapshot, 0, len(snaps))

	for _, s := range snaps {
		shifted := s
		shifted.Timestamp = s.Timestamp.Add(opts.Offset)

		if opts.ClampToNow && shifted.Timestamp.After(now) {
			shifted.Timestamp = now
		}

		out = append(out, shifted)
	}

	return out, nil
}
