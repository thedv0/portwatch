package snapshot

import (
	"testing"
	"time"
)

func dcPort(port int, proto, proc string, pid int) Port {
	return Port{Port: port, Protocol: proto, Process: proc, PID: pid}
}

func makeDCSnap(t time.Time, ports ...Port) Snapshot {
	return Snapshot{Timestamp: t, Ports: ports}
}

func TestBuildDedupeChain_EmptyInput(t *testing.T) {
	result := BuildDedupeChain(nil, DefaultDedupeChainOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestBuildDedupeChain_EntryCount(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeDCSnap(now, dcPort(80, "tcp", "nginx", 1)),
		makeDCSnap(now.Add(time.Minute), dcPort(443, "tcp", "nginx", 2)),
	}
	entries := BuildDedupeChain(snaps, DefaultDedupeChainOptions())
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestBuildDedupeChain_RemovesDuplicates(t *testing.T) {
	now := time.Now()
	snap := makeDCSnap(now,
		dcPort(80, "tcp", "nginx", 1),
		dcPort(80, "tcp", "nginx", 1),
		dcPort(443, "tcp", "nginx", 2),
	)
	entries := BuildDedupeChain([]Snapshot{snap}, DefaultDedupeChainOptions())
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Removed != 1 {
		t.Errorf("expected 1 removed, got %d", entries[0].Removed)
	}
	if len(entries[0].Snapshot.Ports) != 2 {
		t.Errorf("expected 2 ports after dedupe, got %d", len(entries[0].Snapshot.Ports))
	}
}

func TestBuildDedupeChain_NoDuplicates_RemovedIsZero(t *testing.T) {
	now := time.Now()
	snap := makeDCSnap(now,
		dcPort(80, "tcp", "nginx", 1),
		dcPort(443, "tcp", "nginx", 2),
	)
	entries := BuildDedupeChain([]Snapshot{snap}, DefaultDedupeChainOptions())
	if entries[0].Removed != 0 {
		t.Errorf("expected 0 removed, got %d", entries[0].Removed)
	}
}

func TestBuildDedupeChain_TimestampPreserved(t *testing.T) {
	now := time.Now().Round(time.Second)
	snap := makeDCSnap(now, dcPort(22, "tcp", "sshd", 10))
	entries := BuildDedupeChain([]Snapshot{snap}, DefaultDedupeChainOptions())
	if !entries[0].Timestamp.Equal(now) {
		t.Errorf("expected timestamp %v, got %v", now, entries[0].Timestamp)
	}
}
