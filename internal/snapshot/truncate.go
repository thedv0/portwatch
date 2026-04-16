package snapshot

import "time"

// TruncateOptions controls how snapshots are truncated.
type TruncateOptions struct {
	// MaxPorts limits the number of ports retained per snapshot.
	MaxPorts int
	// Before removes ports last seen before this time (zero means no filter).
	Before time.Time
	// Protocols restricts truncation to only these protocols (empty means all).
	Protocols []string
}

// DefaultTruncateOptions returns sensible defaults.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxPorts: 0,
	}
}

// Truncate removes ports from each snapshot according to the given options.
// It returns a new slice of snapshots with modified port lists.
func Truncate(snaps []Snapshot, opts TruncateOptions) []Snapshot {
	protoSet := make(map[string]bool, len(opts.Protocols))
	for _, p := range opts.Protocols {
		protoSet[p] = true
	}

	result := make([]Snapshot, 0, len(snaps))
	for _, s := range snaps {
		filtered := make([]Port, 0, len(s.Ports))
		for _, p := range s.Ports {
			if len(protoSet) > 0 && !protoSet[p.Protocol] {
				filtered = append(filtered, p)
				continue
			}
			if !opts.Before.IsZero() && !p.SeenAt.IsZero() && p.SeenAt.Before(opts.Before) {
				continue
			}
			filtered = append(filtered, p)
		}
		if opts.MaxPorts > 0 && len(filtered) > opts.MaxPorts {
			filtered = filtered[:opts.MaxPorts]
		}
		snap := s
		snap.Ports = filtered
		result = append(result, snap)
	}
	return result
}
