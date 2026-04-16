package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func unPort(proto string, port int, pid int, process string) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, PID: pid, Process: process}
}

func makeUnionSnaps(groups ...[]scanner.Port) []Snapshot {
	snaps := make([]Snapshot, len(groups))
	for i, g := range groups {
		snaps[i] = Snapshot{Ports: g}
	}
	return snaps
}

func TestUnion_EmptyInput(t *testing.T) {
	result, err := Union(nil, DefaultUnionOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty, got %d ports", len(result))
	}
}

func TestUnion_NoDuplicates_ReturnsAll(t *testing.T) {
	snaps := makeUnionSnaps(
		[]scanner.Port{unPort("tcp", 80, 1, "nginx")},
		[]scanner.Port{unPort("tcp", 443, 2, "nginx")},
	)
	result, err := Union(snaps, DefaultUnionOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 ports, got %d", len(result))
	}
}

func TestUnion_DeduplicatesOverlap(t *testing.T) {
	snaps := makeUnionSnaps(
		[]scanner.Port{unPort("tcp", 80, 1, "nginx"), unPort("tcp", 443, 2, "nginx")},
		[]scanner.Port{unPort("tcp", 80, 3, "apache")}, // same proto+port, different pid/process
	)
	result, err := Union(snaps, DefaultUnionOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 unique ports, got %d", len(result))
	}
}

func TestUnion_DedupFalse_KeepsDuplicates(t *testing.T) {
	opts := DefaultUnionOptions()
	opts.Dedup = false
	snaps := makeUnionSnaps(
		[]scanner.Port{unPort("tcp", 80, 1, "nginx")},
		[]scanner.Port{unPort("tcp", 80, 1, "nginx")},
	)
	result, err := Union(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 (duplicates kept), got %d", len(result))
	}
}

func TestUnion_InvalidOptions_ReturnsError(t *testing.T) {
	opts := DefaultUnionOptions()
	opts.KeyFields = nil
	_, err := Union([]Snapshot{{}}, opts)
	if err == nil {
		t.Error("expected error for empty KeyFields")
	}
}

func TestUnion_KeyByPID_DistinguishesSamePort(t *testing.T) {
	opts := DefaultUnionOptions()
	opts.KeyFields = []string{"proto", "port", "pid"}
	snaps := makeUnionSnaps(
		[]scanner.Port{unPort("tcp", 80, 1, "nginx")},
		[]scanner.Port{unPort("tcp", 80, 2, "nginx")}, // same port, different pid
	)
	result, err := Union(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 ports (different PIDs), got %d", len(result))
	}
}
