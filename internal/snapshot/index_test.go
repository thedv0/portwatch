package snapshot

import (
	"testing"
	"time"
)

func idxPort(port int, proto, process string, pid int) PortState {
	return PortState{Port: port, Protocol: proto, Process: process, PID: pid}
}

func makeIdxSnap(ports ...PortState) Snapshot {
	return Snapshot{Timestamp: time.Now(), Ports: ports}
}

func TestIndex_DefaultKeyFields(t *testing.T) {
	snaps := []Snapshot{
		makeIdxSnap(idxPort(80, "tcp", "nginx", 100), idxPort(443, "tcp", "nginx", 100)),
	}
	idx, err := Index(snaps, DefaultIndexOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(idx))
	}
	if _, ok := idx["80/tcp"]; !ok {
		t.Error("expected key 80/tcp")
	}
	if _, ok := idx["443/tcp"]; !ok {
		t.Error("expected key 443/tcp")
	}
}

func TestIndex_MultipleSnaps_Accumulates(t *testing.T) {
	s1 := makeIdxSnap(idxPort(80, "tcp", "nginx", 100))
	s2 := makeIdxSnap(idxPort(80, "tcp", "nginx", 101))
	idx, err := Index([]Snapshot{s1, s2}, DefaultIndexOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := idx["80/tcp"]
	if len(entry.Ports) != 2 {
		t.Fatalf("expected 2 ports in entry, got %d", len(entry.Ports))
	}
}

func TestIndex_ByProcess(t *testing.T) {
	snaps := []Snapshot{
		makeIdxSnap(idxPort(80, "tcp", "nginx", 10), idxPort(8080, "tcp", "nginx", 11)),
	}
	opts := IndexOptions{KeyFields: []string{"process"}}
	idx, err := Index(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry, ok := idx["nginx"]
	if !ok {
		t.Fatal("expected key 'nginx'")
	}
	if len(entry.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(entry.Ports))
	}
}

func TestIndex_EmptyKeyFields_ReturnsError(t *testing.T) {
	_, err := Index([]Snapshot{}, IndexOptions{KeyFields: []string{}})
	if err == nil {
		t.Fatal("expected error for empty key fields")
	}
}

func TestIndex_UnknownField_ReturnsError(t *testing.T) {
	_, err := Index([]Snapshot{}, IndexOptions{KeyFields: []string{"unknown"}})
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestSortedKeys_Order(t *testing.T) {
	idx := map[string]IndexEntry{
		"80/tcp":  {Key: "80/tcp"},
		"22/tcp":  {Key: "22/tcp"},
		"443/tcp": {Key: "443/tcp"},
	}
	keys := SortedKeys(idx)
	if keys[0] != "22/tcp" || keys[1] != "443/tcp" || keys[2] != "80/tcp" {
		t.Errorf("unexpected order: %v", keys)
	}
}
