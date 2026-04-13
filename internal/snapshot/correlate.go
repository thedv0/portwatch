package snapshot

import (
	"fmt"
	"sort"
)

// CorrelateOptions controls how port correlation is performed.
type CorrelateOptions struct {
	// MatchByPID correlates ports that share the same PID across snapshots.
	MatchByPID bool
	// MatchByProcess correlates ports that share the same process name.
	MatchByProcess bool
	// MatchByPort correlates ports that share the same port number and protocol.
	MatchByPort bool
}

// DefaultCorrelateOptions returns sensible defaults.
func DefaultCorrelateOptions() CorrelateOptions {
	return CorrelateOptions{
		MatchByPID:     false,
		MatchByProcess: true,
		MatchByPort:    true,
	}
}

// CorrelatedGroup holds ports from multiple snapshots that are considered
// related under the chosen correlation strategy.
type CorrelatedGroup struct {
	Key    string
	Ports  []Port
	Count  int
}

// Correlate groups ports from a set of snapshots according to opts.
// Each snapshot is a slice of Port values; snaps is the collection of them.
func Correlate(snaps [][]Port, opts CorrelateOptions) []CorrelatedGroup {
	index := make(map[string][]Port)

	for _, snap := range snaps {
		for _, p := range snap {
			keys := correlationKeys(p, opts)
			for _, k := range keys {
				index[k] = append(index[k], p)
			}
		}
	}

	groups := make([]CorrelatedGroup, 0, len(index))
	for k, ports := range index {
		groups = append(groups, CorrelatedGroup{
			Key:   k,
			Ports: ports,
			Count: len(ports),
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Count != groups[j].Count {
			return groups[i].Count > groups[j].Count
		}
		return groups[i].Key < groups[j].Key
	})

	return groups
}

func correlationKeys(p Port, opts CorrelateOptions) []string {
	seen := make(map[string]struct{})
	var keys []string

	add := func(k string) {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}

	if opts.MatchByPort {
		add(fmt.Sprintf("port:%s:%d", p.Protocol, p.Port))
	}
	if opts.MatchByProcess && p.Process != "" {
		add(fmt.Sprintf("process:%s", p.Process))
	}
	if opts.MatchByPID && p.PID > 0 {
		add(fmt.Sprintf("pid:%d", p.PID))
	}

	return keys
}
