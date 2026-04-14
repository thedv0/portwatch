package snapshot

import (
	"math"
	"testing"
	"time"
)

func makeSmoothSnap(t time.Time, count int) Snapshot {
	ports := make([]Port, count)
	for i := range ports {
		ports[i] = Port{Port: i + 1, Protocol: "tcp"}
	}
	return Snapshot{Timestamp: t, Ports: ports}
}

func TestSmooth_RequiresMinSnapshots(t *testing.T) {
	opts := DefaultSmoothOptions()
	_, err := Smooth([]Snapshot{makeSmoothSnap(time.Now(), 3)}, opts)
	if err == nil {
		t.Fatal("expected error for too few snapshots")
	}
}

func TestSmooth_InvalidAlpha(t *testing.T) {
	opts := DefaultSmoothOptions()
	opts.Alpha = 0
	snaps := []Snapshot{
		makeSmoothSnap(time.Now(), 2),
		makeSmoothSnap(time.Now().Add(time.Second), 4),
	}
	_, err := Smooth(snaps, opts)
	if err == nil {
		t.Fatal("expected error for alpha=0")
	}
}

func TestSmooth_PointCount(t *testing.T) {
	base := time.Now()
	snaps := []Snapshot{
		makeSmoothSnap(base, 2),
		makeSmoothSnap(base.Add(time.Second), 4),
		makeSmoothSnap(base.Add(2*time.Second), 6),
	}
	pts, err := Smooth(snaps, DefaultSmoothOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pts) != 3 {
		t.Fatalf("expected 3 points, got %d", len(pts))
	}
}

func TestSmooth_FirstPointEqualsRaw(t *testing.T) {
	base := time.Now()
	snaps := []Snapshot{
		makeSmoothSnap(base, 5),
		makeSmoothSnap(base.Add(time.Second), 10),
	}
	pts, err := Smooth(snaps, DefaultSmoothOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pts[0].Smoothed != 5.0 {
		t.Errorf("expected first smoothed=5, got %v", pts[0].Smoothed)
	}
}

func TestSmooth_EMAFormula(t *testing.T) {
	base := time.Now()
	opts := SmoothOptions{Alpha: 0.5, MinSnapshots: 2}
	snaps := []Snapshot{
		makeSmoothSnap(base, 4),
		makeSmoothSnap(base.Add(time.Second), 8),
	}
	pts, err := Smooth(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// EMA[1] = 0.5*8 + 0.5*4 = 6
	want := 6.0
	if math.Abs(pts[1].Smoothed-want) > 1e-9 {
		t.Errorf("expected smoothed=%v, got %v", want, pts[1].Smoothed)
	}
}

func TestSmooth_RawCountPreserved(t *testing.T) {
	base := time.Now()
	snaps := []Snapshot{
		makeSmoothSnap(base, 3),
		makeSmoothSnap(base.Add(time.Second), 7),
	}
	pts, _ := Smooth(snaps, DefaultSmoothOptions())
	if pts[1].RawCount != 7 {
		t.Errorf("expected raw count 7, got %d", pts[1].RawCount)
	}
}
