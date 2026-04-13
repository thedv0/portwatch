package snapshot

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func ruPort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, PID: 100, Process: "test"}
}

func makeRollupSnaps(offsets []time.Duration, ports [][]scanner.Port) []Snapshot {
	now := time.Now()
	var snaps []Snapshot
	for i, off := range offsets {
		snaps = append(snaps, Snapshot{
			Timestamp: now.Add(-off),
			Ports:     ports[i],
		})
	}
	return snaps
}

func TestRollup_EmptyInput(t *testing.T) {
	res := Rollup(nil, DefaultRollupOptions())
	if len(res.Ports) != 0 {
		t.Fatalf("expected 0 ports, got %d", len(res.Ports))
	}
	if res.SnapshotCount != 0 {
		t.Errorf("expected 0 snapshots included")
	}
}

func TestRollup_MergesAllPorts(t *testing.T) {
	snaps := makeRollupSnaps(
		[]time.Duration{10 * time.Second, 20 * time.Second},
		[][]scanner.Port{
			{ruPort(80, "tcp"), ruPort(443, "tcp")},
			{ruPort(8080, "tcp")},
		},
	)
	res := Rollup(snaps, DefaultRollupOptions())
	if len(res.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(res.Ports))
	}
	if res.SnapshotCount != 2 {
		t.Errorf("expected 2 snapshots included, got %d", res.SnapshotCount)
	}
}

func TestRollup_WindowFiltersOldSnaps(t *testing.T) {
	snaps := makeRollupSnaps(
		[]time.Duration{5 * time.Second, 2 * time.Minute},
		[][]scanner.Port{
			{ruPort(80, "tcp")},
			{ruPort(9999, "tcp")},
		},
	)
	opts := DefaultRollupOptions()
	opts.Window = 30 * time.Second
	res := Rollup(snaps, opts)
	if res.SnapshotCount != 1 {
		t.Errorf("expected 1 snapshot in window, got %d", res.SnapshotCount)
	}
	if len(res.Ports) != 1 || res.Ports[0].Port != 80 {
		t.Errorf("expected only port 80, got %v", res.Ports)
	}
}

func TestRollup_MinOccurrencesFilters(t *testing.T) {
	snaps := makeRollupSnaps(
		[]time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second},
		[][]scanner.Port{
			{ruPort(80, "tcp"), ruPort(22, "tcp")},
			{ruPort(80, "tcp")},
			{ruPort(80, "tcp")},
		},
	)
	opts := DefaultRollupOptions()
	opts.MinOccurrences = 3
	res := Rollup(snaps, opts)
	if len(res.Ports) != 1 || res.Ports[0].Port != 80 {
		t.Errorf("expected only port 80 with 3 occurrences, got %v", res.Ports)
	}
	if res.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", res.Dropped)
	}
}

func TestRollup_DeduplicatesWithUniqueOnly(t *testing.T) {
	snaps := makeRollupSnaps(
		[]time.Duration{1 * time.Second},
		[][]scanner.Port{
			{ruPort(80, "tcp"), ruPort(80, "tcp"), ruPort(443, "tcp")},
		},
	)
	opts := DefaultRollupOptions()
	opts.UniqueOnly = true
	res := Rollup(snaps, opts)
	if len(res.Ports) > 2 {
		t.Errorf("expected at most 2 unique ports, got %d", len(res.Ports))
	}
}
