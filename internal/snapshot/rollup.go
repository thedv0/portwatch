package snapshot

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// RollupOptions controls how snapshots are rolled up into a single summary.
type RollupOptions struct {
	// Window is the time range to include. Zero means include all.
	Window time.Duration
	// UniqueOnly keeps only unique port+protocol combinations.
	UniqueOnly bool
	// MinOccurrences filters ports that appear fewer times than this threshold.
	MinOccurrences int
}

// DefaultRollupOptions returns sensible defaults.
func DefaultRollupOptions() RollupOptions {
	return RollupOptions{
		Window:         0,
		UniqueOnly:     false,
		MinOccurrences: 1,
	}
}

// RollupResult holds the output of a rollup operation.
type RollupResult struct {
	Timestamp  time.Time
	Ports      []scanner.Port
	SnapshotCount int
	Dropped    int
}

// Rollup merges multiple snapshots into a single deduplicated result,
// applying optional windowing and occurrence filtering.
func Rollup(snaps []Snapshot, opts RollupOptions) RollupResult {
	now := time.Now()
	counts := make(map[string]int)
	merged := make(map[string]scanner.Port)

	included := 0
	for _, s := range snaps {
		if opts.Window > 0 && now.Sub(s.Timestamp) > opts.Window {
			continue
		}
		included++
		for _, p := range s.Ports {
			k := rollupKey(p)
			counts[k]++
			if _, seen := merged[k]; !seen {
				merged[k] = p
			}
		}
	}

	dropped := 0
	var result []scanner.Port
	for k, p := range merged {
		if counts[k] < opts.MinOccurrences {
			dropped++
			continue
		}
		result = append(result, p)
	}

	if opts.UniqueOnly {
		seen := make(map[string]bool)
		var uniq []scanner.Port
		for _, p := range result {
			k := rollupKey(p)
			if !seen[k] {
				seen[k] = true
				uniq = append(uniq, p)
			}
		}
		result = uniq
	}

	return RollupResult{
		Timestamp:     now,
		Ports:         result,
		SnapshotCount: included,
		Dropped:       dropped,
	}
}

func rollupKey(p scanner.Port) string {
	return itoa(p.Port) + "/" + p.Protocol
}
