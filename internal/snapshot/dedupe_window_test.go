package snapshot

import (
	"testing"
	"time"
)

func dwPort(proto string, port int) Port {
	return Port{Proto: proto, Port: port, Process: "proc"}
}

func makeDWSnap(ts time.Time, ports ...Port) Snapshot {
	return Snapshot{Timestamp: ts, Ports: ports}
}

var dwBase = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestDedupeWindow_EmptyInput(t *testing.T) {
	result, err := DedupeWindow(nil, DefaultDedupeWindowOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Snapshots) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(result.Snapshots))
	}
}

func TestDedupeWindow_InvalidWindowSize(t *testing.T) {
	opts := DefaultDedupeWindowOptions()
	opts.WindowSize = 0
	_, err := DedupeWindow([]Snapshot{makeDWSnap(dwBase, dwPort("tcp", 80))}, opts)
	if err == nil {
		t.Fatal("expected error for zero WindowSize")
	}
}

func TestDedupeWindow_FirstSnapshotUnchanged(t *testing.T) {
	snaps := []Snapshot{
		makeDWSnap(dwBase, dwPort("tcp", 80), dwPort("tcp", 443)),
	}
	result, err := DedupeWindow(snaps, DefaultDedupeWindowOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Snapshots[0].Ports) != 2 {
		t.Errorf("expected 2 ports in first snapshot, got %d", len(result.Snapshots[0].Ports))
	}
}

func TestDedupeWindow_DropsRepeatedPortsInWindow(t *testing.T) {
	opts := DefaultDedupeWindowOptions()
	opts.WindowSize = 10 * time.Minute
	snaps := []Snapshot{
		makeDWSnap(dwBase, dwPort("tcp", 80), dwPort("tcp", 443)),
		makeDWSnap(dwBase.Add(2*time.Minute), dwPort("tcp", 80), dwPort("tcp", 8080)),
	}
	result, err := DedupeWindow(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// port 80 should be dropped from second snap; 8080 is new
	if len(result.Snapshots[1].Ports) != 1 {
		t.Errorf("expected 1 port in second snapshot, got %d", len(result.Snapshots[1].Ports))
	}
	if result.Snapshots[1].Ports[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", result.Snapshots[1].Ports[0].Port)
	}
}

func TestDedupeWindow_OutsideWindowRetainsAll(t *testing.T) {
	opts := DefaultDedupeWindowOptions()
	opts.WindowSize = 1 * time.Minute
	snaps := []Snapshot{
		makeDWSnap(dwBase, dwPort("tcp", 80)),
		makeDWSnap(dwBase.Add(10*time.Minute), dwPort("tcp", 80)),
	}
	result, err := DedupeWindow(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Snapshots[1].Ports) != 1 {
		t.Errorf("expected port 80 retained outside window, got %d ports", len(result.Snapshots[1].Ports))
	}
}

func TestDedupeWindow_DroppedCount(t *testing.T) {
	opts := DefaultDedupeWindowOptions()
	opts.WindowSize = 10 * time.Minute
	snaps := []Snapshot{
		makeDWSnap(dwBase, dwPort("tcp", 80), dwPort("tcp", 443)),
		makeDWSnap(dwBase.Add(1*time.Minute), dwPort("tcp", 80), dwPort("tcp", 443)),
	}
	result, err := DedupeWindow(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Dropped != 2 {
		t.Errorf("expected 2 dropped, got %d", result.Dropped)
	}
}
