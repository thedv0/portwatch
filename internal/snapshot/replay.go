package snapshot

import (
	"errors"
	"sort"
	"time"
)

// ReplayOptions controls how a sequence of snapshots is replayed.
type ReplayOptions struct {
	// StartAt filters snapshots taken at or after this time. Zero means no lower bound.
	StartAt time.Time
	// EndAt filters snapshots taken at or before this time. Zero means no upper bound.
	EndAt time.Time
	// MaxFrames limits the number of frames returned. 0 means unlimited.
	MaxFrames int
	// Reverse replays from newest to oldest when true.
	Reverse bool
}

// DefaultReplayOptions returns sensible defaults.
func DefaultReplayOptions() ReplayOptions {
	return ReplayOptions{
		MaxFrames: 0,
		Reverse:   false,
	}
}

// ReplayFrame is a single step in a replay sequence.
type ReplayFrame struct {
	Index     int
	Timestamp time.Time
	Ports     []PortState
	Diff      DiffResult
}

// Replay takes an ordered slice of snapshots and produces a sequence of
// ReplayFrames, each containing the diff relative to the previous frame.
func Replay(snaps []Snapshot, opts ReplayOptions) ([]ReplayFrame, error) {
	if len(snaps) == 0 {
		return nil, nil
	}

	filtered := filterByTime(snaps, opts.StartAt, opts.EndAt)
	if len(filtered) == 0 {
		return nil, nil
	}

	sort.Slice(filtered, func(i, j int) bool {
		if opts.Reverse {
			return filtered[i].Timestamp.After(filtered[j].Timestamp)
		}
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	if opts.MaxFrames > 0 && len(filtered) > opts.MaxFrames {
		filtered = filtered[:opts.MaxFrames]
	}

	frames := make([]ReplayFrame, 0, len(filtered))
	var prev []PortState
	for i, snap := range filtered {
		if snap.Ports == nil {
			return nil, errors.New("replay: snapshot contains nil ports slice")
		}
		d := Diff(prev, snap.Ports)
		frames = append(frames, ReplayFrame{
			Index:     i,
			Timestamp: snap.Timestamp,
			Ports:     snap.Ports,
			Diff:      d,
		})
		prev = snap.Ports
	}
	return frames, nil
}

func filterByTime(snaps []Snapshot, start, end time.Time) []Snapshot {
	out := snaps[:0:0]
	for _, s := range snaps {
		if !start.IsZero() && s.Timestamp.Before(start) {
			continue
		}
		if !end.IsZero() && s.Timestamp.After(end) {
			continue
		}
		out = append(out, s)
	}
	return out
}
