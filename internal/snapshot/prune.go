package snapshot

import (
	"sort"
	"time"
)

// PruneOptions controls how ports are pruned from a snapshot.
type PruneOptions struct {
	// MaxAge removes ports whose last-seen time is older than this duration.
	// Zero means no age-based pruning.
	MaxAge time.Duration

	// AllowedProtocols, if non-empty, removes ports not matching any listed protocol.
	AllowedProtocols []string

	// PortBlacklist removes ports whose port number appears in this list.
	PortBlacklist []int

	// MaxPorts, if > 0, retains only the first N ports sorted by port number.
	MaxPorts int
}

// DefaultPruneOptions returns a PruneOptions with no restrictions.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{}
}

// Prune filters ports from the given slice according to PruneOptions.
// It returns a new slice; the original is not modified.
func Prune(ports []PortState, opts PruneOptions) []PortState {
	blacklist := make(map[int]struct{}, len(opts.PortBlacklist))
	for _, p := range opts.PortBlacklist {
		blacklist[p] = struct{}{}
	}

	allowed := make(map[string]struct{}, len(opts.AllowedProtocols))
	for _, proto := range opts.AllowedProtocols {
		allowed[proto] = struct{}{}
	}

	now := time.Now()
	result := make([]PortState, 0, len(ports))

	for _, p := range ports {
		if _, blocked := blacklist[p.Port]; blocked {
			continue
		}
		if len(allowed) > 0 {
			if _, ok := allowed[p.Protocol]; !ok {
				continue
			}
		}
		if opts.MaxAge > 0 && !p.SeenAt.IsZero() {
			if now.Sub(p.SeenAt) > opts.MaxAge {
				continue
			}
		}
		result = append(result, p)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})

	if opts.MaxPorts > 0 && len(result) > opts.MaxPorts {
		result = result[:opts.MaxPorts]
	}

	return result
}
