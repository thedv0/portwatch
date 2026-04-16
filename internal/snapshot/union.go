package snapshot

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// DefaultUnionOptions returns a UnionOptions with sensible defaults.
func DefaultUnionOptions() UnionOptions {
	return UnionOptions{
		KeyFields: []string{"proto", "port"},
		Dedup:     true,
	}
}

// UnionOptions controls how snapshots are combined.
type UnionOptions struct {
	// KeyFields determines how ports are deduplicated ("proto", "port", "pid", "process").
	KeyFields []string
	// Dedup removes duplicate ports across all snapshots when true.
	Dedup bool
}

// Validate returns an error if the options are invalid.
func (o UnionOptions) Validate() error {
	if len(o.KeyFields) == 0 {
		return fmt.Errorf("union: KeyFields must not be empty")
	}
	return nil
}

// Union combines all ports from every snapshot into a single slice.
// When Dedup is true, duplicate ports (by key) are removed, keeping the
// first occurrence in iteration order.
func Union(snaps []Snapshot, opts UnionOptions) ([]scanner.Port, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	if len(snaps) == 0 {
		return []scanner.Port{}, nil
	}

	seen := make(map[string]struct{})
	var result []scanner.Port

	for _, snap := range snaps {
		for _, p := range snap.Ports {
			if !opts.Dedup {
				result = append(result, p)
				continue
			}
			k := unionKey(p, opts.KeyFields)
			if _, exists := seen[k]; !exists {
				seen[k] = struct{}{}
				result = append(result, p)
			}
		}
	}

	if result == nil {
		return []scanner.Port{}, nil
	}
	return result, nil
}

func unionKey(p scanner.Port, fields []string) string {
	key := ""
	for _, f := range fields {
		switch f {
		case "proto":
			key += p.Protocol + "|"
		case "port":
			key += fmt.Sprintf("%d|", p.Port)
		case "pid":
			key += fmt.Sprintf("%d|", p.PID)
		case "process":
			key += p.Process + "|"
		}
	}
	return key
}
