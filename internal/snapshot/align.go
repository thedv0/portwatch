package snapshot

import (
	"errors"
	"sort"
	"time"
)

// AlignOptions controls how snapshots are aligned to a common time grid.
type AlignOptions struct {
	// BucketSize is the duration of each time bucket.
	BucketSize time.Duration
	// Tolerance is the maximum distance a snapshot can be from a bucket
	// boundary and still be considered aligned to it.
	Tolerance time.Duration
}

// DefaultAlignOptions returns sensible defaults for AlignOptions.
func DefaultAlignOptions() AlignOptions {
	return AlignOptions{
		BucketSize: time.Minute,
		Tolerance:  10 * time.Second,
	}
}

// AlignedBucket represents a time bucket with the snapshots aligned to it.
type AlignedBucket struct {
	BucketTime time.Time
	Snapshots  []Snapshot
}

// Align groups snapshots into fixed-size time buckets. Each snapshot is
// placed into the bucket whose boundary it is nearest to, provided it falls
// within the configured tolerance. Snapshots outside any bucket's tolerance
// are dropped. Buckets are returned in ascending time order.
func Align(snaps []Snapshot, opts AlignOptions) ([]AlignedBucket, error) {
	if opts.BucketSize <= 0 {
		return nil, errors.New("align: BucketSize must be positive")
	}
	if opts.Tolerance < 0 {
		return nil, errors.New("align: Tolerance must be non-negative")
	}
	if len(snaps) == 0 {
		return []AlignedBucket{}, nil
	}

	buckets := map[int64]*AlignedBucket{}

	for _, s := range snaps {
		// Round to nearest bucket boundary.
		boundary := s.Timestamp.Round(opts.BucketSize)
		dist := s.Timestamp.Sub(boundary)
		if dist < 0 {
			dist = -dist
		}
		if dist > opts.Tolerance {
			continue
		}
		key := boundary.UnixNano()
		if _, ok := buckets[key]; !ok {
			buckets[key] = &AlignedBucket{BucketTime: boundary}
		}
		buckets[key].Snapshots = append(buckets[key].Snapshots, s)
	}

	result := make([]AlignedBucket, 0, len(buckets))
	for _, b := range buckets {
		result = append(result, *b)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].BucketTime.Before(result[j].BucketTime)
	})
	return result, nil
}
