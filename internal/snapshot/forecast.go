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
	if opts.Horizon <= 0 {
		return ForecastResult{}, errors.New("forecast: horizon must be positive")
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

// AtTime returns the forecasted port count at a specific future time by
// evaluating the regression line. The provided time must be after the last
// snapshot used to build the result, otherwise an error is returned.
func (r ForecastResult) AtTime(t time.Time) (float64, error) {
	if r.GeneratedAt.IsZero() {
		return 0, errors.New("forecast: result is empty")
	}
	if !t.After(r.GeneratedAt) {
		return 0, errors.New("forecast: requested time must be after forecast generation time")
	}
	x := t.Sub(r.GeneratedAt).Seconds()
	predicted := r.Slope*x + r.Intercept
	if predicted < 0 {
		predicted = 0
	}
	return math.Round(predicted*100) / 100, nil
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
