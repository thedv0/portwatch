package snapshot

import (
	"fmt"
	"sort"

	"github.com/netwatch/portwatch/internal/scanner"
)

// ClusterOptions controls how ports are clustered.
type ClusterOptions struct {
	// By is the field to cluster on: "port", "process", "protocol", "pid".
	By string
	// MinSize discards clusters with fewer than MinSize members.
	MinSize int
}

// DefaultClusterOptions returns sensible defaults.
func DefaultClusterOptions() ClusterOptions {
	return ClusterOptions{
		By:      "process",
		MinSize: 1,
	}
}

// Cluster groups ports from multiple snapshots into named clusters.
// Each cluster key is derived from the chosen field.
func Cluster(snaps []scanner.Snapshot, opts ClusterOptions) map[string][]scanner.Port {
	if opts.By == "" {
		opts.By = DefaultClusterOptions().By
	}
	if opts.MinSize < 1 {
		opts.MinSize = 1
	}

	acc := make(map[string][]scanner.Port)

	for _, s := range snaps {
		for _, p := range s.Ports {
			key := clusterKey(p, opts.By)
			acc[key] = append(acc[key], p)
		}
	}

	// Remove clusters below MinSize.
	for k, v := range acc {
		if len(v) < opts.MinSize {
			delete(acc, k)
		}
	}

	// Sort each cluster by port for determinism.
	for k := range acc {
		sort.Slice(acc[k], func(i, j int) bool {
			return acc[k][i].Port < acc[k][j].Port
		})
	}

	return acc
}

func clusterKey(p scanner.Port, by string) string {
	switch by {
	 case "port":
		return fmt.Sprintf("%d", p.Port)
	case "protocol":
		if p.Protocol == "" {
			return "unknown"
		}
		return p.Protocol
	case "pid":
		return fmt.Sprintf("%d", p.PID)
	default: // "process"
		if p.Process == "" {
			return "unknown"
		}
		return p.Process
	}
}
