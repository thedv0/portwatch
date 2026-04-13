package snapshot

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PartitionOptions controls how snapshots are partitioned into time buckets.
type PartitionOptions struct {
	// BucketDuration defines the width of each time bucket.
	BucketDuration time.Duration
	// MaxBuckets caps the number of returned buckets (0 = unlimited).
	MaxBuckets int
}

// DefaultPartitionOptions returns sensible defaults.
func DefaultPartitionOptions() PartitionOptions {
	return PartitionOptions{
		BucketDuration: time.Hour,
		MaxBuckets:     0,
	}
}

// Bucket holds all ports observed within a single time window.
type Bucket struct {
	Key   string           // formatted start time of the bucket
	Start time.Time
	End   time.Time
	Ports []scanner.Port
}

// Partition groups a slice of snapshots into time-aligned buckets.
// Snapshots are expected to carry a Timestamp field; those with zero
// timestamps are placed in a special "unknown" bucket.
func Partition(snaps []Snapshot, opts PartitionOptions) []Bucket {
	if len(snaps) == 0 {
		return nil
	}
	if opts.BucketDuration <= 0 {
		opts.BucketDuration = DefaultPartitionOptions().BucketDuration
	}

	index := make(map[string]*Bucket)
	order := []string{}

	for _, s := range snaps {
		var key string
		var start, end time.Time

		if s.Timestamp.IsZero() {
			key = "unknown"
		} else {
			trunc := s.Timestamp.Truncate(opts.BucketDuration)
			start = trunc
			end = trunc.Add(opts.BucketDuration)
			key = fmt.Sprintf("%d", trunc.UnixNano())
		}

		if _, ok := index[key]; !ok {
			index[key] = &Bucket{Key: key, Start: start, End: end}
			order = append(order, key)
		}
		index[key].Ports = append(index[key].Ports, s.Ports...)
	}

	result := make([]Bucket, 0, len(order))
	for _, k := range order {
		result = append(result, *index[k])
		if opts.MaxBuckets > 0 && len(result) >= opts.MaxBuckets {
			break
		}
	}
	return result
}
