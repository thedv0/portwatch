package snapshot

import (
	"errors"
	"sort"
	"time"
)

// ResampleOptions controls how port samples are resampled into buckets.
type ResampleOptions struct {
	// BucketSize is the duration of each bucket.
	BucketSize time.Duration
	// Aggregator determines how ports within a bucket are combined.
	// Supported values: "union" (default), "intersect".
	Aggregator string
}

// DefaultResampleOptions returns sensible defaults.
func DefaultResampleOptions() ResampleOptions {
	return ResampleOptions{
		BucketSize: time.Minute,
		Aggregator: "union",
	}
}

// Bucket holds the resampled ports for a time window.
type Bucket struct {
	Start time.Time
	End   time.Time
	Ports []Port
}

// Resample groups the provided snapshots into fixed-size time buckets and
// merges the ports within each bucket according to the chosen aggregator.
func Resample(snaps []Snapshot, opts ResampleOptions) ([]Bucket, error) {
	if opts.BucketSize <= 0 {
		return nil, errors.New("resample: BucketSize must be positive")
	}
	if opts.Aggregator == "" {
		opts.Aggregator = "union"
	}
	if opts.Aggregator != "union" && opts.Aggregator != "intersect" {
		return nil, errors.New("resample: unknown aggregator: " + opts.Aggregator)
	}
	if len(snaps) == 0 {
		return nil, nil
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].Timestamp.Before(snaps[j].Timestamp)
	})

	origin := snaps[0].Timestamp.Truncate(opts.BucketSize)
	bucketMap := make(map[int64][]Snapshot)
	for _, s := range snaps {
		offset := int64(s.Timestamp.Sub(origin) / opts.BucketSize)
		bucketMap[offset] = append(bucketMap[offset], s)
	}

	keys := make([]int64, 0, len(bucketMap))
	for k := range bucketMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	mergeOpts := DefaultMergeOptions()
	mergeOpts.Strategy = opts.Aggregator

	buckets := make([]Bucket, 0, len(keys))
	for _, k := range keys {
		group := bucketMap[k]
		start := origin.Add(time.Duration(k) * opts.BucketSize)
		end := start.Add(opts.BucketSize)
		var merged []Port
		for _, s := range group {
			merged = Merge(merged, s.Ports, mergeOpts)
		}
		buckets = append(buckets, Bucket{Start: start, End: end, Ports: merged})
	}
	return buckets, nil
}
