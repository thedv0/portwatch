package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// IndexOptions controls how ports are indexed.
type IndexOptions struct {
	// KeyFields determines which fields form the index key.
	// Valid values: "port", "protocol", "pid", "process"
	KeyFields []string
}

// DefaultIndexOptions returns sensible defaults.
func DefaultIndexOptions() IndexOptions {
	return IndexOptions{
		KeyFields: []string{"port", "protocol"},
	}
}

// IndexEntry holds all ports that share the same composite key.
type IndexEntry struct {
	Key   string
	Ports []PortState
}

// Index builds a map from composite key to matching ports across one or more snapshots.
func Index(snaps []Snapshot, opts IndexOptions) (map[string]IndexEntry, error) {
	if len(opts.KeyFields) == 0 {
		return nil, fmt.Errorf("index: at least one key field required")
	}
	for _, f := range opts.KeyFields {
		switch f {
		case "port", "protocol", "pid", "process":
		default:
			return nil, fmt.Errorf("index: unknown key field %q", f)
		}
	}

	result := make(map[string]IndexEntry)
	for _, snap := range snaps {
		for _, p := range snap.Ports {
			k := buildIndexKey(p, opts.KeyFields)
			e := result[k]
			e.Key = k
			e.Ports = append(e.Ports, p)
			result[k] = e
		}
	}
	return result, nil
}

// SortedKeys returns the index keys in lexicographic order.
func SortedKeys(idx map[string]IndexEntry) []string {
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func buildIndexKey(p PortState, fields []string) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "port":
			parts = append(parts, fmt.Sprintf("%d", p.Port))
		case "protocol":
			parts = append(parts, strings.ToLower(p.Protocol))
		case "pid":
			parts = append(parts, fmt.Sprintf("%d", p.PID))
		case "process":
			parts = append(parts, p.Process)
		}
	}
	return strings.Join(parts, "/")
}
