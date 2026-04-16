package snapshot

import (
	"fmt"
	"time"
)

// DedupeWindowOptions configures the DedupeWindow operation.
type DedupeWindowOptions struct {
	// WindowSize is the duration within which duplicate ports are collapsed.
	WindowSize time.Duration
	// KeyFields controls what fields define a duplicate.
	KeyFields []string
}

// DefaultDedupeWindowOptions returns sensible defaults.
func DefaultDedupeWindowOptions() DedupeWindowOptions {
	return DedupeWindowOptions{
		WindowSize: 5 * time.Minute,
		KeyFields:  []string{"proto", "port"},
	}
}

// DedupeWindowResult holds the output of a DedupeWindow operation.
type DedupeWindowResult struct {
	Snapshots []Snapshot
	Dropped   int
	WindowSize time.Duration
}

// DedupeWindow collapses snapshots within a sliding time window, removing
// ports that appear unchanged across consecutive snapshots in the window.
func DedupeWindow(snaps []Snapshot, opts DedupeWindowOptions) (DedupeWindowResult, error) {
	if opts.WindowSize <= 0 {
		return DedupeWindowResult{}, fmt.Errorf("dedupe_window: WindowSize must be positive")
	}
	if len(snaps) == 0 {
		return DedupeWindowResult{WindowSize: opts.WindowSize}, nil
	}

	result := make([]Snapshot, 0, len(snaps))
	dropped := 0

	for i, snap := range snaps {
		if i == 0 {
			result = append(result, snap)
			continue
		}
		prev := result[len(result)-1]
		if snap.Timestamp.Sub(prev.Timestamp) <= opts.WindowSize {
			unique := dedupeWindowFilter(prev.Ports, snap.Ports, opts.KeyFields)
			dropped += len(snap.Ports) - len(unique)
			snap.Ports = unique
		}
		result = append(result, snap)
	}

	return DedupeWindowResult{
		Snapshots:  result,
		Dropped:    dropped,
		WindowSize: opts.WindowSize,
	}, nil
}

func dedupeWindowFilter(prev, curr []Port, keyFields []string) []Port {
	seen := make(map[string]struct{}, len(prev))
	for _, p := range prev {
		seen[dedupeWindowKey(p, keyFields)] = struct{}{}
	}
	out := make([]Port, 0, len(curr))
	for _, p := range curr {
		if _, exists := seen[dedupeWindowKey(p, keyFields)]; !exists {
			out = append(out, p)
		}
	}
	return out
}

func dedupeWindowKey(p Port, fields []string) string {
	key := ""
	for _, f := range fields {
		switch f {
		case "proto":
			key += p.Proto + "|"
		case "port":
			key += itoa(p.Port) + "|"
		case "pid":
			key += itoa(p.PID) + "|"
		case "process":
			key += p.Process + "|"
		}
	}
	return key
}
