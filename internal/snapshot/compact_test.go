package snapshot

import (
	"testing"
	"time"
)

func csnap(ts time.Time, ports ...Port) Snapshot {
	return Snapshot{Timestamp: ts, Ports: ports}
}

func csport(port int, proto string) Port {
	return Port{Port: port, Protocol: proto, PID: 1, Process: "test"}
}

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestCompact_EmptyInput(t *testing.T) {
	r := Compact(nil, DefaultCompactOptions())
	if r.Before != 0 || r.After != 0 {
		t.Fatalf("expected zero counts, got before=%d after=%d", r.Before, r.After)
	}
}

func TestCompact_NoMerge_NoDropWhenUnderLimit(t *testing.T) {
	snaps := []Snapshot{
		csnap(base, csport(80, "tcp")),
		csnap(base.Add(10*time.Minute), csport(443, "tcp")),
	}
	opts := DefaultCompactOptions()
	opts.MaxSnapshots = 10
	opts.MergeWindow = 0
	r := Compact(snaps, opts)
	if r.Before != 2 || r.After != 2 {
		t.Fatalf("expected 2 snapshots, got before=%d after=%d", r.Before, r.After)
	}
	if r.Dropped != 0 {
		t.Fatalf("expected 0 dropped, got %d", r.Dropped)
	}
}

func TestCompact_MergesWithinWindow(t *testing.T) {
	snaps := []Snapshot{
		csnap(base, csport(80, "tcp")),
		csnap(base.Add(30*time.Second), csport(443, "tcp")),
		csnap(base.Add(45*time.Second), csport(8080, "tcp")),
	}
	opts := DefaultCompactOptions()
	opts.MergeWindow = time.Minute
	opts.MaxSnapshots = 0
	r := Compact(snaps, opts)
	if r.After != 1 {
		t.Fatalf("expected 1 merged snapshot, got %d", r.After)
	}
	if r.Merged != 2 {
		t.Fatalf("expected merged=2, got %d", r.Merged)
	}
	if len(r.Snapshots[0].Ports) != 3 {
		t.Fatalf("expected 3 ports in merged snapshot, got %d", len(r.Snapshots[0].Ports))
	}
}

func TestCompact_MergeDeduplicatesPorts(t *testing.T) {
	p := csport(80, "tcp")
	snaps := []Snapshot{
		csnap(base, p),
		csnap(base.Add(10*time.Second), p),
	}
	opts := DefaultCompactOptions()
	opts.MergeWindow = time.Minute
	opts.MaxSnapshots = 0
	r := Compact(snaps, opts)
	if len(r.Snapshots[0].Ports) != 1 {
		t.Fatalf("expected 1 deduplicated port, got %d", len(r.Snapshots[0].Ports))
	}
}

func TestCompact_DropsOldWhenOverLimit(t *testing.T) {
	var snaps []Snapshot
	for i := 0; i < 10; i++ {
		snaps = append(snaps, csnap(base.Add(time.Duration(i)*time.Hour), csport(80+i, "tcp")))
	}
	opts := DefaultCompactOptions()
	opts.MaxSnapshots = 5
	opts.MergeWindow = 0
	r := Compact(snaps, opts)
	if r.After != 5 {
		t.Fatalf("expected 5 snapshots after compaction, got %d", r.After)
	}
	if r.Dropped != 5 {
		t.Fatalf("expected 5 dropped, got %d", r.Dropped)
	}
}

func TestCompact_MergedSnapshotUsesLatestTimestamp(t *testing.T) {
	later := base.Add(50 * time.Second)
	snaps := []Snapshot{
		csnap(base, csport(80, "tcp")),
		csnap(later, csport(443, "tcp")),
	}
	opts := DefaultCompactOptions()
	opts.MergeWindow = time.Minute
	opts.MaxSnapshots = 0
	r := Compact(snaps, opts)
	if !r.Snapshots[0].Timestamp.Equal(later) {
		t.Fatalf("expected timestamp %v, got %v", later, r.Snapshots[0].Timestamp)
	}
}
