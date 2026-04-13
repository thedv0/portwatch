package snapshot

import (
	"errors"
	"math"
	"time"
)

// ForecastOptions controls how port count forecasting behaves.
type ForecastOptions struct {
	// Horizon is how far into the future to forecast.
	Horizon time.Duration
	// Steps is the number of forecast points to generate.
	Steps int
}

// ForecastPoint represents a single predicted data point.
type ForecastPoint struct {
	At        time.Time
	PortCount float64
}

// ForecastResult holds the output of a forecast run.
type ForecastResult struct {
	GeneratedAt time.Time
	Horizon     time.Duration
	Points      []ForecastPoint
	Slope       float64
	Intercept   float64
}

// DefaultForecastOptions returns sensible defaults.
func DefaultForecastOptions() ForecastOptions {
	return ForecastOptions{
		Horizon: 1 * time.Hour,
		Steps:   6,
	}
}

// Forecast predicts future open port counts using linear regression over
// historical snapshots. At least two snapshots are required.
func Forecast(snaps []Snapshot, opts ForecastOptions) (ForecastResult, error) {
	if len(snaps) < 2 {
		return ForecastResult{}, errors.New("forecast: at least two snapshots required")
	}
	if opts.Steps <= 0 {
		return ForecastResult{}, errors.New("forecast: steps must be positive")
	}

	// Build (x, y) pairs where x = seconds since first snapshot, y = port count.
	origin := snaps[0].Timestamp
	xs := make([]float64, len(snaps))
	ys := make([]float64, len(snaps))
	for i, s := range snaps {
		xs[i] = s.Timestamp.Sub(origin).Seconds()
		ys[i] = float64(len(s.Ports))
	}

	slope, intercept := linearRegression(xs, ys)

	// Generate forecast points evenly spaced over the horizon.
	stepDur := opts.Horizon / time.Duration(opts.Steps)
	last := snaps[len(snaps)-1].Timestamp
	points := make([]ForecastPoint, opts.Steps)
	for i := 0; i < opts.Steps; i++ {
		at := last.Add(stepDur * time.Duration(i+1))
		x := at.Sub(origin).Seconds()
		predicted := slope*x + intercept
		if predicted < 0 {
			predicted = 0
		}
		points[i] = ForecastPoint{At: at, PortCount: math.Round(predicted*100) / 100}
	}

	return ForecastResult{
		GeneratedAt: time.Now(),
		Horizon:     opts.Horizon,
		Points:      points,
		Slope:       slope,
		Intercept:   intercept,
	}, nil
}

func linearRegression(xs, ys []float64) (slope, intercept float64) {
	n := float64(len(xs))
	var sumX, sumY, sumXY, sumX2 float64
	for i := range xs {
		sumX += xs[i]
		sumY += ys[i]
		sumXY += xs[i] * ys[i]
		sumX2 += xs[i] * xs[i]
	}
	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return 0, sumY / n
	}
	slope = (n*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / n
	return slope, intercept
}
