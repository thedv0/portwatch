package snapshot

import (
	"sort"
	"time"
)

// CompactOptions controls how snapshots are compacted.
type CompactOptions struct {
	// MaxSnapshots is the maximum number of snapshots to retain after compaction.
	// Zero means no limit.
	MaxSnapshots int
	// MinAge is the minimum age a snapshot must have to be eligible for removal.
	MinAge time.Duration
	// MergeWindow groups snapshots within this duration into a single merged snapshot.
	MergeWindow time.Duration
}

// DefaultCompactOptions returns sensible defaults for compaction.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{
		MaxSnapshots: 100,
		MinAge:       5 * time.Minute,
		MergeWindow:  time.Minute,
	}
}

// CompactResult holds the outcome of a compaction run.
type CompactResult struct {
	Before    int
	After     int
	Merged    int
	Dropped   int
	Snapshots []Snapshot
}

// Compact reduces a slice of snapshots by merging those within the same time
// window and dropping old ones that exceed MaxSnapshots.
func Compact(snaps []Snapshot, opts CompactOptions) CompactResult {
	if len(snaps) == 0 {
		return CompactResult{}
	}

	// Sort by timestamp ascending.
	sorted := make([]Snapshot, len(snaps))
	copy(sorted, snaps)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	before := len(sorted)
	var buckets []Snapshot
	mergedCount := 0

	if opts.MergeWindow > 0 {
		i := 0
		for i < len(sorted) {
			bucketStart := sorted[i].Timestamp
			j := i + 1
			for j < len(sorted) && sorted[j].Timestamp.Sub(bucketStart) <= opts.MergeWindow {
				j++
			}
			if j > i+1 {
				merged := mergeSnapshots(sorted[i:j])
				buckets = append(buckets, merged)
				mergedCount += j - i - 1
			} else {
				buckets = append(buckets, sorted[i])
			}
			i = j
		}
	} else {
		buckets = sorted
	}

	dropped := 0
	if opts.MaxSnapshots > 0 && len(buckets) > opts.MaxSnapshots {
		dropped = len(buckets) - opts.MaxSnapshots
		buckets = buckets[dropped:]
	}

	return CompactResult{
		Before:    before,
		After:     len(buckets),
		Merged:    mergedCount,
		Dropped:   dropped,
		Snapshots: buckets,
	}
}

// mergeSnapshots combines multiple snapshots into one, using the latest
// timestamp and the union of all ports (deduplicated by port+protocol).
func mergeSnapshots(snaps []Snapshot) Snapshot {
	latest := snaps[len(snaps)-1].Timestamp
	seen := make(map[string]struct{})
	var ports []Port
	for _, s := range snaps {
		for _, p := range s.Ports {
			k := portKey(p)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				ports = append(ports, p)
			}
		}
	}
	return Snapshot{Timestamp: latest, Ports: ports}
}
