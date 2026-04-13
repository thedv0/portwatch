package snapshot

import (
	"testing"
	"time"
)

func makeResampleSnap(ts time.Time, ports []Port) Snapshot {
	return Snapshot{Timestamp: ts, Ports: ports}
}

func rsPort(proto string, port int) Port {
	return Port{Protocol: proto, Port: port, PID: 1, Process: "p"}
}

func TestResample_EmptyInput(t *testing.T) {
	buckets, err := Resample(nil, DefaultResampleOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(buckets))
	}
}

func TestResample_InvalidBucketSize(t *testing.T) {
	opts := DefaultResampleOptions()
	opts.BucketSize = 0
	_, err := Resample([]Snapshot{{}}, opts)
	if err == nil {
		t.Fatal("expected error for zero BucketSize")
	}
}

func TestResample_UnknownAggregator(t *testing.T) {
	opts := DefaultResampleOptions()
	opts.Aggregator = "median"
	_, err := Resample([]Snapshot{{}}, opts)
	if err == nil {
		t.Fatal("expected error for unknown aggregator")
	}
}

func TestResample_SingleBucket(t *testing.T) {
	now := time.Now().Truncate(time.Minute)
	snaps := []Snapshot{
		makeResampleSnap(now.Add(5*time.Second), []Port{rsPort("tcp", 80)}),
		makeResampleSnap(now.Add(30*time.Second), []Port{rsPort("tcp", 443)}),
	}
	opts := DefaultResampleOptions()
	buckets, err := Resample(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if len(buckets[0].Ports) != 2 {
		t.Errorf("expected 2 ports in bucket, got %d", len(buckets[0].Ports))
	}
}

func TestResample_MultipleBuckets(t *testing.T) {
	now := time.Now().Truncate(time.Minute)
	snaps := []Snapshot{
		makeResampleSnap(now, []Port{rsPort("tcp", 80)}),
		makeResampleSnap(now.Add(2*time.Minute), []Port{rsPort("tcp", 443)}),
	}
	opts := DefaultResampleOptions()
	buckets, err := Resample(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
}

func TestResample_BucketTimeRange(t *testing.T) {
	now := time.Now().Truncate(time.Minute)
	snaps := []Snapshot{
		makeResampleSnap(now.Add(10*time.Second), []Port{rsPort("tcp", 22)}),
	}
	opts := DefaultResampleOptions()
	buckets, err := Resample(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket")
	}
	b := buckets[0]
	if !b.End.Equal(b.Start.Add(opts.BucketSize)) {
		t.Errorf("bucket end should be start + BucketSize")
	}
}

func TestResample_IntersectAggregator(t *testing.T) {
	now := time.Now().Truncate(time.Minute)
	snaps := []Snapshot{
		makeResampleSnap(now.Add(5*time.Second), []Port{rsPort("tcp", 80), rsPort("tcp", 443)}),
		makeResampleSnap(now.Add(20*time.Second), []Port{rsPort("tcp", 80)}),
	}
	opts := DefaultResampleOptions()
	opts.Aggregator = "intersect"
	buckets, err := Resample(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket")
	}
	if len(buckets[0].Ports) != 1 {
		t.Errorf("intersect: expected 1 common port, got %d", len(buckets[0].Ports))
	}
}
