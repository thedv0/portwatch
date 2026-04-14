package snapshot

import (
	"fmt"
	"time"
)

// SmoothOptions controls the exponential moving average smoothing of port counts.
type SmoothOptions struct {
	// Alpha is the smoothing factor in (0, 1]. Higher values weight recent data more.
	Alpha float64
	// MinSnapshots is the minimum number of snapshots required to produce output.
	MinSnapshots int
}

// DefaultSmoothOptions returns sensible defaults for smoothing.
func DefaultSmoothOptions() SmoothOptions {
	return SmoothOptions{
		Alpha:        0.3,
		MinSnapshots: 2,
	}
}

// SmoothedPoint represents a single EMA-smoothed observation.
type SmoothedPoint struct {
	Timestamp  time.Time
	RawCount   int
	Smoothed   float64
}

// Smooth applies exponential moving average smoothing to the open-port counts
// across a series of snapshots. Snapshots must be ordered oldest-first.
func Smooth(snaps []Snapshot, opts SmoothOptions) ([]SmoothedPoint, error) {
	if opts.Alpha <= 0 || opts.Alpha > 1 {
		return nil, fmt.Errorf("smooth: alpha must be in (0, 1], got %v", opts.Alpha)
	}
	if len(snaps) < opts.MinSnapshots {
		return nil, fmt.Errorf("smooth: need at least %d snapshots, got %d", opts.MinSnapshots, len(snaps))
	}

	points := make([]SmoothedPoint, 0, len(snaps))
	var ema float64

	for i, s := range snaps {
		raw := float64(len(s.Ports))
		if i == 0 {
			ema = raw
		} else {
			ema = opts.Alpha*raw + (1-opts.Alpha)*ema
		}
		points = append(points, SmoothedPoint{
			Timestamp: s.Timestamp,
			RawCount:  len(s.Ports),
			Smoothed:  ema,
		})
	}
	return points, nil
}
