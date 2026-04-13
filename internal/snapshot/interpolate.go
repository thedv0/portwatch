package snapshot

import (
	"errors"
	"sort"
	"time"
)

// InterpolateOptions controls how missing snapshots are filled in.
type InterpolateOptions struct {
	// MaxGap is the maximum duration between snapshots to fill.
	// Gaps larger than this are left as-is.
	MaxGap time.Duration
	// Method is the interpolation method: "forward" or "zero".
	Method string
}

// DefaultInterpolateOptions returns sensible defaults.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		MaxGap: 5 * time.Minute,
		Method: "forward",
	}
}

// InterpolateResult holds the output of an interpolation pass.
type InterpolateResult struct {
	Snapshots  []Snapshot
	Filled     int
	Timestamp  time.Time
}

// Interpolate fills gaps between snapshots using the chosen method.
// Snapshots are sorted by timestamp before processing.
func Interpolate(snaps []Snapshot, step time.Duration, opts InterpolateOptions) (InterpolateResult, error) {
	if step <= 0 {
		return InterpolateResult{}, errors.New("interpolate: step must be positive")
	}
	if len(snaps) < 2 {
		return InterpolateResult{Snapshots: snaps, Timestamp: time.Now()}, nil
	}

	sorted := make([]Snapshot, len(snaps))
	copy(sorted, snaps)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	var out []Snapshot
	filled := 0

	for i := 0; i < len(sorted); i++ {
		out = append(out, sorted[i])
		if i+1 >= len(sorted) {
			break
		}
		gap := sorted[i+1].Timestamp.Sub(sorted[i].Timestamp)
		if gap <= step || gap > opts.MaxGap {
			continue
		}
		cursor := sorted[i].Timestamp.Add(step)
		for cursor.Before(sorted[i+1].Timestamp) {
			var synthetic Snapshot
			synthetic.Timestamp = cursor
			if opts.Method == "forward" {
				synthetic.Ports = sorted[i].Ports
			}
			out = append(out, synthetic)
			filled++
			cursor = cursor.Add(step)
		}
	}

	return InterpolateResult{
		Snapshots: out,
		Filled:    filled,
		Timestamp: time.Now(),
	}, nil
}
