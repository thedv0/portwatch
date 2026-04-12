package snapshot

import (
	"testing"
)

func pts(pairs ...int64) []TrendPoint {
	var out []TrendPoint
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, TrendPoint{Timestamp: pairs[i], Value: float64(pairs[i+1])})
	}
	return out
}

func TestAnalyzeTrend_Stable(t *testing.T) {
	points := pts(1, 5, 2, 5, 3, 5, 4, 5)
	res := AnalyzeTrend(points, DefaultTrendOptions())
	if res.Direction != TrendStable {
		t.Fatalf("expected stable, got %s", res.Direction)
	}
}

func TestAnalyzeTrend_Up(t *testing.T) {
	points := pts(1, 1, 2, 3, 3, 5, 4, 7)
	res := AnalyzeTrend(points, DefaultTrendOptions())
	if res.Direction != TrendUp {
		t.Fatalf("expected up, got %s", res.Direction)
	}
	if res.Slope <= 0 {
		t.Fatalf("expected positive slope, got %f", res.Slope)
	}
}

func TestAnalyzeTrend_Down(t *testing.T) {
	points := pts(1, 10, 2, 7, 3, 4, 4, 1)
	res := AnalyzeTrend(points, DefaultTrendOptions())
	if res.Direction != TrendDown {
		t.Fatalf("expected down, got %s", res.Direction)
	}
	if res.Slope >= 0 {
		t.Fatalf("expected negative slope, got %f", res.Slope)
	}
}

func TestAnalyzeTrend_SinglePoint(t *testing.T) {
	points := []TrendPoint{{Timestamp: 1, Value: 42}}
	res := AnalyzeTrend(points, DefaultTrendOptions())
	if res.Direction != TrendStable {
		t.Fatalf("single point should be stable, got %s", res.Direction)
	}
}

func TestAnalyzeTrend_Empty(t *testing.T) {
	res := AnalyzeTrend(nil, DefaultTrendOptions())
	if res.Direction != TrendStable {
		t.Fatalf("empty should be stable, got %s", res.Direction)
	}
}

func TestAnalyzeTrend_SortsPoints(t *testing.T) {
	// Provide points out of order; slope should still be computed correctly.
	points := []TrendPoint{
		{Timestamp: 4, Value: 8},
		{Timestamp: 1, Value: 2},
		{Timestamp: 3, Value: 6},
		{Timestamp: 2, Value: 4},
	}
	res := AnalyzeTrend(points, DefaultTrendOptions())
	if res.Direction != TrendUp {
		t.Fatalf("expected up after sorting, got %s", res.Direction)
	}
	for i := 1; i < len(res.Points); i++ {
		if res.Points[i].Timestamp < res.Points[i-1].Timestamp {
			t.Fatal("points not sorted by timestamp")
		}
	}
}

func TestDefaultTrendOptions(t *testing.T) {
	opts := DefaultTrendOptions()
	if opts.SlopeThreshold <= 0 {
		t.Fatal("default slope threshold should be positive")
	}
}
