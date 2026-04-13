package snapshot

import (
	"testing"
	"time"
)

func makeForecastSnap(t time.Time, portCount int) Snapshot {
	ports := make([]Port, portCount)
	for i := range ports {
		ports[i] = Port{Port: i + 1, Protocol: "tcp", PID: 100 + i}
	}
	return Snapshot{Timestamp: t, Ports: ports}
}

func TestForecast_RequiresTwoSnapshots(t *testing.T) {
	snaps := []Snapshot{makeForecastSnap(time.Now(), 5)}
	_, err := Forecast(snaps, DefaultForecastOptions())
	if err == nil {
		t.Fatal("expected error for single snapshot")
	}
}

func TestForecast_StepsMustBePositive(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeForecastSnap(now, 5),
		makeForecastSnap(now.Add(time.Minute), 6),
	}
	opts := DefaultForecastOptions()
	opts.Steps = 0
	_, err := Forecast(snaps, opts)
	if err == nil {
		t.Fatal("expected error for zero steps")
	}
}

func TestForecast_PointCount(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeForecastSnap(now, 4),
		makeForecastSnap(now.Add(time.Minute), 6),
		makeForecastSnap(now.Add(2*time.Minute), 8),
	}
	opts := DefaultForecastOptions()
	opts.Steps = 4
	res, err := Forecast(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Points) != 4 {
		t.Errorf("expected 4 points, got %d", len(res.Points))
	}
}

func TestForecast_UpwardTrend(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeForecastSnap(now, 2),
		makeForecastSnap(now.Add(time.Minute), 4),
		makeForecastSnap(now.Add(2*time.Minute), 6),
	}
	res, err := Forecast(snaps, DefaultForecastOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Slope <= 0 {
		t.Errorf("expected positive slope for upward trend, got %f", res.Slope)
	}
	for _, p := range res.Points {
		if p.PortCount < 0 {
			t.Errorf("forecast point should not be negative: %f", p.PortCount)
		}
	}
}

func TestForecast_PointsAreChronological(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeForecastSnap(now, 5),
		makeForecastSnap(now.Add(time.Minute), 7),
	}
	res, err := Forecast(snaps, DefaultForecastOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(res.Points); i++ {
		if !res.Points[i].At.After(res.Points[i-1].At) {
			t.Errorf("points not in chronological order at index %d", i)
		}
	}
}

func TestForecast_DefaultOptions(t *testing.T) {
	opts := DefaultForecastOptions()
	if opts.Steps <= 0 {
		t.Errorf("default steps should be positive, got %d", opts.Steps)
	}
	if opts.Horizon <= 0 {
		t.Errorf("default horizon should be positive")
	}
}
