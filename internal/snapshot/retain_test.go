package snapshot

import (
	"testing"
	"time"
)

func makeRetainSnap(ago time.Duration, ports []Port) Snapshot {
	return Snapshot{
		Timestamp: time.Now().Add(-ago),
		Ports:     ports,
	}
}

func TestRetain_DefaultOptions_Valid(t *testing.T) {
	opts := DefaultRetainOptions()
	if err := opts.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRetain_KeepsRecentSnapshots(t *testing.T) {
	snaps := []Snapshot{
		makeRetainSnap(1*time.Hour, nil),
		makeRetainSnap(2*time.Hour, nil),
	}
	opts := DefaultRetainOptions() // MaxAge = 72h
	result, err := Retain(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(result))
	}
}

func TestRetain_RemovesOldSnapshots(t *testing.T) {
	snaps := []Snapshot{
		makeRetainSnap(1*time.Hour, nil),
		makeRetainSnap(100*time.Hour, nil),
	}
	opts := RetainOptions{MaxAge: 48 * time.Hour, MinCount: 0, MaxCount: 0}
	result, err := Retain(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(result))
	}
}

func TestRetain_MinCountPreservesOld(t *testing.T) {
	snaps := []Snapshot{
		makeRetainSnap(200*time.Hour, nil),
	}
	opts := RetainOptions{MaxAge: 48 * time.Hour, MinCount: 1, MaxCount: 0}
	result, err := Retain(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("MinCount should preserve at least 1 snapshot, got %d", len(result))
	}
}

func TestRetain_MaxCountLimitsResults(t *testing.T) {
	snaps := []Snapshot{
		makeRetainSnap(1*time.Hour, nil),
		makeRetainSnap(2*time.Hour, nil),
		makeRetainSnap(3*time.Hour, nil),
	}
	opts := RetainOptions{MaxAge: 0, MinCount: 0, MaxCount: 2}
	result, err := Retain(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(result))
	}
}

func TestRetain_InvalidMaxAge_ReturnsError(t *testing.T) {
	opts := RetainOptions{MaxAge: -1 * time.Second}
	_, err := Retain(nil, opts)
	if err == nil {
		t.Fatal("expected error for negative MaxAge")
	}
}

func TestRetain_SortedNewestFirst(t *testing.T) {
	old := makeRetainSnap(10*time.Hour, nil)
	new_ := makeRetainSnap(1*time.Hour, nil)
	snaps := []Snapshot{old, new_}
	opts := RetainOptions{MaxAge: 0, MinCount: 0, MaxCount: 0}
	result, err := Retain(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) < 2 {
		t.Fatal("expected at least 2 results")
	}
	if !result[0].Timestamp.After(result[1].Timestamp) {
		t.Error("expected results sorted newest first")
	}
}
