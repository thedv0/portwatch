package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// ProjectOptions controls how future port states are projected.
type ProjectOptions struct {
	// Steps is the number of future time steps to project.
	Steps int
	// StepSize is the duration between projected steps.
	StepSize time.Duration
	// MinSnapshots is the minimum number of snapshots required.
	MinSnapshots int
}

// DefaultProjectOptions returns sensible defaults.
func DefaultProjectOptions() ProjectOptions {
	return ProjectOptions{
		Steps:        5,
		StepSize:     time.Minute,
		MinSnapshots: 2,
	}
}

// ProjectedPoint represents a single projected port-count value at a future time.
type ProjectedPoint struct {
	At    time.Time
	Count float64
}

// ProjectResult holds the output of a projection run.
type ProjectResult struct {
	Points    []ProjectedPoint
	BaseCount float64
	Slope     float64
	GeneratedAt time.Time
}

// Project uses linear extrapolation over historical snapshots to estimate
// future open-port counts.
func Project(snaps []PortSnapshot, opts ProjectOptions) (ProjectResult, error) {
	if opts.Steps <= 0 {
		return ProjectResult{}, fmt.Errorf("project: steps must be positive")
	}
	if opts.StepSize <= 0 {
		return ProjectResult{}, fmt.Errorf("project: step size must be positive")
	}
	if len(snaps) < opts.MinSnapshots {
		return ProjectResult{}, fmt.Errorf("project: need at least %d snapshots, got %d", opts.MinSnapshots, len(snaps))
	}

	sorted := make([]PortSnapshot, len(snaps))
	copy(sorted, snaps)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	n := float64(len(sorted))
	var sumX, sumY, sumXY, sumX2 float64
	origin := sorted[0].Timestamp
	for i, s := range sorted {
		x := float64(i)
		y := float64(len(s.Ports))
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denom := n*sumX2 - sumX*sumX
	var slope, intercept float64
	if denom != 0 {
		slope = (n*sumXY - sumX*sumY) / denom
		intercept = (sumY - slope*sumX) / n
	}

	last := sorted[len(sorted)-1]
	points := make([]ProjectedPoint, opts.Steps)
	for i := 0; i < opts.Steps; i++ {
		xFut := n + float64(i)
		at := last.Timestamp.Add(time.Duration(i+1) * opts.StepSize)
		_ = origin
		points[i] = ProjectedPoint{
			At:    at,
			Count: intercept + slope*xFut,
		}
	}

	return ProjectResult{
		Points:      points,
		BaseCount:   float64(len(last.Ports)),
		Slope:       slope,
		GeneratedAt: time.Now(),
	}, nil
}
