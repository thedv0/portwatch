package snapshot

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func wsnap(ts time.Time, ports ...scanner.Port) Snapshot {
	return Snapshot{Timestamp: ts, Ports: ports}
}

func wport(p int) scanner.Port { return scanner.Port{Port: p, Protocol: "tcp"} }

var (
	base   = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	minus2 = base.Add(-2 * time.Hour)
	minus1 = base.Add(-1 * time.Hour)
	plus1  = base.Add(1 * time.Hour)
	plus2  = base.Add(2 * time.Hour)
)

func TestWindow_ReturnsSnapsInRange(t *testing.T) {
	snaps := []Snapshot{
		wsnap(minus2, wport(80)),
		wsnap(minus1, wport(443)),
		wsnap(base, wport(8080)),
		wsnap(plus1, wport(9090)),
		wsnap(plus2, wport(22)),
	}
	opts := WindowOptions{Start: minus1, End: plus1}
	res, err := Window(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 {
		t.Errorf("expected Total=3, got %d", res.Total)
	}
	if len(res.Snapshots) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(res.Snapshots))
	}
}

func TestWindow_OrderedByTimestamp(t *testing.T) {
	snaps := []Snapshot{
		wsnap(plus1, wport(9090)),
		wsnap(minus1, wport(443)),
		wsnap(base, wport(8080)),
	}
	opts := WindowOptions{Start: minus1, End: plus1}
	res, err := Window(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(res.Snapshots); i++ {
		if res.Snapshots[i].Timestamp.Before(res.Snapshots[i-1].Timestamp) {
			t.Errorf("snapshots not sorted at index %d", i)
		}
	}
}

func TestWindow_MaxSnapshotsCaps(t *testing.T) {
	snaps := []Snapshot{
		wsnap(minus1, wport(1)),
		wsnap(base, wport(2)),
		wsnap(plus1, wport(3)),
	}
	opts := WindowOptions{Start: minus1, End: plus1, MaxSnapshots: 2}
	res, err := Window(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 {
		t.Errorf("expected Total=3, got %d", res.Total)
	}
	if len(res.Snapshots) != 2 {
		t.Errorf("expected 2 snapshots after cap, got %d", len(res.Snapshots))
	}
}

func TestWindow_EmptyInput(t *testing.T) {
	opts := WindowOptions{Start: minus1, End: plus1}
	res, err := Window(nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 0 || len(res.Snapshots) != 0 {
		t.Error("expected empty result for nil input")
	}
}

func TestWindow_InvalidOptions_EndBeforeStart(t *testing.T) {
	opts := WindowOptions{Start: plus1, End: minus1}
	_, err := Window(nil, opts)
	if err == nil {
		t.Error("expected error for End before Start")
	}
}

func TestDefaultWindowOptions_Valid(t *testing.T) {
	opts := DefaultWindowOptions()
	if err := opts.Validate(); err != nil {
		t.Errorf("default options should be valid: %v", err)
	}
}
