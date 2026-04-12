package snapshot

import (
	"math"
	"sort"
)

// TrendDirection indicates whether a metric is increasing, decreasing, or stable.
type TrendDirection string

const (
	TrendUp     TrendDirection = "up"
	TrendDown   TrendDirection = "down"
	TrendStable TrendDirection = "stable"
)

// TrendPoint represents a single observation over time.
type TrendPoint struct {
	Timestamp int64
	Value     float64
}

// TrendResult holds the computed trend for a metric.
type TrendResult struct {
	Direction TrendDirection
	Slope     float64
	Points    []TrendPoint
}

// TrendOptions controls trend analysis behaviour.
type TrendOptions struct {
	// SlopeThreshold is the minimum absolute slope to be considered non-stable.
	SlopeThreshold float64
}

// DefaultTrendOptions returns sensible defaults.
func DefaultTrendOptions() TrendOptions {
	return TrendOptions{SlopeThreshold: 0.1}
}

// AnalyzeTrend computes a linear regression slope over the given points and
// returns a TrendResult describing the direction of change.
func AnalyzeTrend(points []TrendPoint, opts TrendOptions) TrendResult {
	sorted := make([]TrendPoint, len(points))
	copy(sorted, points)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp < sorted[j].Timestamp
	})

	result := TrendResult{Points: sorted, Direction: TrendStable}
	if len(sorted) < 2 {
		return result
	}

	slope := linearSlope(sorted)
	result.Slope = slope
	switch {
	case slope > opts.SlopeThreshold:
		result.Direction = TrendUp
	case slope < -opts.SlopeThreshold:
		result.Direction = TrendDown
	default:
		result.Direction = TrendStable
	}
	return result
}

// linearSlope computes the slope of the best-fit line through the points.
func linearSlope(pts []TrendPoint) float64 {
	n := float64(len(pts))
	var sumX, sumY, sumXY, sumX2 float64
	for _, p := range pts {
		x := float64(p.Timestamp)
		y := p.Value
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denom := n*sumX2 - sumX*sumX
	if math.Abs(denom) < 1e-9 {
		return 0
	}
	return (n*sumXY - sumX*sumY) / denom
}
