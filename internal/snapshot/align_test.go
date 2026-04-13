package snapshot

import (
	"testing"
	"time"
)

func alignSnap(ts time.Time, ports ...Port) Snapshot {
	return Snapshot{Timestamp: ts, Ports: ports}
}

func TestAlign_EmptyInput(t *testing.T) {
	buckets, err := Align(nil, DefaultAlignOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(buckets))
	}
}

func TestAlign_InvalidBucketSize(t *testing.T) {
	opts := DefaultAlignOptions()
	opts.BucketSize = 0
	_, err := Align([]Snapshot{}, opts)
	if err == nil {
		t.Fatal("expected error for zero BucketSize")
	}
}

func TestAlign_NegativeTolerance(t *testing.T) {
	opts := DefaultAlignOptions()
	opts.Tolerance = -1
	_, err := Align([]Snapshot{}, opts)
	if err == nil {
		t.Fatal("expected error for negative Tolerance")
	}
}

func TestAlign_SnapsGroupedIntoBuckets(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	opts := AlignOptions{BucketSize: time.Minute, Tolerance: 10 * time.Second}

	snaps := []Snapshot{
		alignSnap(base.Add(2 * time.Second)),           // rounds to 12:00
		alignSnap(base.Add(58 * time.Second)),          // rounds to 12:01
		alignSnap(base.Add(time.Minute + 3*time.Second)), // rounds to 12:01
	}

	buckets, err := Align(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if len(buckets[0].Snapshots) != 1 {
		t.Errorf("bucket 0: expected 1 snap, got %d", len(buckets[0].Snapshots))
	}
	if len(buckets[1].Snapshots) != 2 {
		t.Errorf("bucket 1: expected 2 snaps, got %d", len(buckets[1].Snapshots))
	}
}

func TestAlign_DropsSnapsOutsideTolerance(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 30, 0, time.UTC) // exactly 30s from boundary
	opts := AlignOptions{BucketSize: time.Minute, Tolerance: 10 * time.Second}

	snaps := []Snapshot{alignSnap(base)}
	buckets, err := Align(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 0 {
		t.Errorf("expected 0 buckets (snap out of tolerance), got %d", len(buckets))
	}
}

func TestAlign_BucketsAreOrdered(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	opts := AlignOptions{BucketSize: time.Minute, Tolerance: 10 * time.Second}

	snaps := []Snapshot{
		alignSnap(base.Add(3 * time.Minute)),
		alignSnap(base.Add(1 * time.Minute)),
		alignSnap(base.Add(2 * time.Minute)),
	}

	buckets, err := Align(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(buckets); i++ {
		if !buckets[i-1].BucketTime.Before(buckets[i].BucketTime) {
			t.Errorf("buckets not ordered at index %d", i)
		}
	}
}
