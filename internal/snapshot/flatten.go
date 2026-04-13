package snapshot

import (
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// FlattenOptions controls how multiple snapshots are flattened into one.
type FlattenOptions struct {
	// Deduplicate removes duplicate ports across snapshots.
	Deduplicate bool
	// SortByPort sorts the result by port number ascending.
	SortByPort bool
	// IncludeProtocols restricts output to the given protocols (empty = all).
	IncludeProtocols []string
}

// DefaultFlattenOptions returns sensible defaults.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Deduplicate: true,
		SortByPort:  true,
	}
}

// Flatten merges a slice of port snapshots into a single deduplicated list.
func Flatten(snaps [][]scanner.Port, opts FlattenOptions) []scanner.Port {
	protoSet := make(map[string]struct{}, len(opts.IncludeProtocols))
	for _, p := range opts.IncludeProtocols {
		protoSet[p] = struct{}{}
	}

	seen := make(map[string]struct{})
	var result []scanner.Port

	for _, snap := range snaps {
		for _, p := range snap {
			if len(protoSet) > 0 {
				if _, ok := protoSet[p.Protocol]; !ok {
					continue
				}
			}
			if opts.Deduplicate {
				key := flattenKey(p)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
			}
			result = append(result, p)
		}
	}

	if opts.SortByPort {
		sort.Slice(result, func(i, j int) bool {
			if result[i].Port != result[j].Port {
				return result[i].Port < result[j].Port
			}
			return result[i].Protocol < result[j].Protocol
		})
	}

	return result
}

func flattenKey(p scanner.Port) string {
	return p.Protocol + ":" + itoa(p.Port)
}
