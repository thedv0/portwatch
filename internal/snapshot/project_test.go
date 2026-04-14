package snapshot

import (
	"testing"
	"time"
)

func makeProjectSnap(t time.Time, count int) PortSnapshot {
	ports := make([]Port, count)
	for i := range ports {
		ports[i] = Port{Port: i + 1, Protocol: "tcp"}
	}
	return PortSnapshot{Timestamp: t, Ports: ports}
}

func TestProject_RequiresTwoSnapshots(t *testing.T) {
	opts := DefaultProjectOptions()
	_, err := Project([]PortSnapshot{makeProjectSnap(time.Now(), 3)}, opts)
	if err == nil {
		t.Fatal("expected error for single snapshot")
	}
}

func TestProject_StepsMustBePositive(t *testing.T) {
	opts := DefaultProjectOptions()
	opts.Steps = 0
	snaps := []PortSnapshot{
		makeProjectSnap(time.Now(), 2),
		makeProjectSnap(time.Now().Add(time.Minute), 4),
	}
	_, err := Project(snaps, opts)
	if err == nil {
		t.Fatal("expected error for zero steps")
	}
}

func TestProject_PointCount(t *testing.T) {
	base := time.Now()
	snaps := []PortSnapshot{
		makeProjectSnap(base, 2),
		makeProjectSnap(base.Add(time.Minute), 4),
		makeProjectSnap(base.Add(2*time.Minute), 6),
	}
	opts := DefaultProjectOptions()
	opts.Steps = 4
	res, err := Project(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Points) != 4 {
		t.Fatalf("expected 4 points, got %d", len(res.Points))
	}
}

func TestProject_UpwardTrend(t *testing.T) {
	base := time.Now()
	snaps := []PortSnapshot{
		makeProjectSnap(base, 2),
		makeProjectSnap(base.Add(time.Minute), 4),
		makeProjectSnap(base.Add(2*time.Minute), 6),
	}
	opts := DefaultProjectOptions()
	opts.Steps = 3
	res, err := Project(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Slope <= 0 {
		t.Errorf("expected positive slope for upward trend, got %f", res.Slope)
	}
}

func TestProject_PointsAreChronological(t *testing.T) {
	base := time.Now()
	snaps := []PortSnapshot{
		makeProjectSnap(base, 3),
		makeProjectSnap(base.Add(time.Minute), 5),
	}
	opts := DefaultProjectOptions()
	opts.Steps = 3
	res, err := Project(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(res.Points); i++ {
		if !res.Points[i].At.After(res.Points[i-1].At) {
			t.Errorf("points not chronological at index %d", i)
		}
	}
}

func TestProject_BaseCountMatchesLastSnapshot(t *testing.T) {
	base := time.Now()
	snaps := []PortSnapshot{
		makeProjectSnap(base, 3),
		makeProjectSnap(base.Add(time.Minute), 7),
	}
	opts := DefaultProjectOptions()
	res, err := Project(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BaseCount != 7 {
		t.Errorf("expected base count 7, got %f", res.BaseCount)
	}
}
