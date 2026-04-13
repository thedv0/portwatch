package snapshot

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func psnap(ts time.Time, ports ...scanner.Port) Snapshot {
	return Snapshot{Timestamp: ts, Ports: ports}
}

func pport(num int) scanner.Port {
	return scanner.Port{Port: num, Protocol: "tcp"}
}

var baseTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestPartition_EmptyInput(t *testing.T) {
	result := Partition(nil, DefaultPartitionOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestPartition_SingleBucket(t *testing.T) {
	snaps := []Snapshot{
		psnap(baseTime, pport(80)),
		psnap(baseTime.Add(10*time.Minute), pport(443)),
	}
	opts := DefaultPartitionOptions() // 1-hour buckets
	buckets := Partition(snaps, opts)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if len(buckets[0].Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(buckets[0].Ports))
	}
}

func TestPartition_MultipleBuckets(t *testing.T) {
	snaps := []Snapshot{
		psnap(baseTime, pport(80)),
		psnap(baseTime.Add(2*time.Hour), pport(443)),
		psnap(baseTime.Add(4*time.Hour), pport(8080)),
	}
	opts := DefaultPartitionOptions()
	buckets := Partition(snaps, opts)
	if len(buckets) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(buckets))
	}
}

func TestPartition_MaxBuckets(t *testing.T) {
	snaps := []Snapshot{
		psnap(baseTime, pport(80)),
		psnap(baseTime.Add(2*time.Hour), pport(443)),
		psnap(baseTime.Add(4*time.Hour), pport(8080)),
	}
	opts := DefaultPartitionOptions()
	opts.MaxBuckets = 2
	buckets := Partition(snaps, opts)
	if len(buckets) != 2 {
		t.Errorf("expected 2 buckets due to MaxBuckets, got %d", len(buckets))
	}
}

func TestPartition_ZeroTimestamp_UnknownBucket(t *testing.T) {
	snaps := []Snapshot{
		{Ports: []scanner.Port{pport(22)}},
		{Ports: []scanner.Port{pport(3306)}},
	}
	buckets := Partition(snaps, DefaultPartitionOptions())
	if len(buckets) != 1 {
		t.Fatalf("expected 1 unknown bucket, got %d", len(buckets))
	}
	if buckets[0].Key != "unknown" {
		t.Errorf("expected key 'unknown', got %q", buckets[0].Key)
	}
	if len(buckets[0].Ports) != 2 {
		t.Errorf("expected 2 ports in unknown bucket, got %d", len(buckets[0].Ports))
	}
}

func TestPartition_BucketStartEnd(t *testing.T) {
	snaps := []Snapshot{psnap(baseTime, pport(80))}
	opts := DefaultPartitionOptions()
	buckets := Partition(snaps, opts)
	if len(buckets) != 1 {
		t.Fatal("expected 1 bucket")
	}
	expectedStart := baseTime.Truncate(opts.BucketDuration)
	expectedEnd := expectedStart.Add(opts.BucketDuration)
	if !buckets[0].Start.Equal(expectedStart) {
		t.Errorf("start mismatch: got %v want %v", buckets[0].Start, expectedStart)
	}
	if !buckets[0].End.Equal(expectedEnd) {
		t.Errorf("end mismatch: got %v want %v", buckets[0].End, expectedEnd)
	}
}
