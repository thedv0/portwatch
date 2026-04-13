package snapshot

import (
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// PivotField selects the dimension to pivot on.
type PivotField string

const (
	PivotByProtocol PivotField = "protocol"
	PivotByProcess  PivotField = "process"
	PivotByPort     PivotField = "port"
)

// PivotOptions controls Pivot behaviour.
type PivotOptions struct {
	Field   PivotField
	SortKeys bool
}

// DefaultPivotOptions returns sensible defaults.
func DefaultPivotOptions() PivotOptions {
	return PivotOptions{
		Field:    PivotByProtocol,
		SortKeys: true,
	}
}

// PivotResult maps a pivot key to the set of ports observed under that key
// across all provided snapshots.
type PivotResult struct {
	Field   PivotField
	Buckets map[string][]scanner.Port
	Keys    []string // ordered when SortKeys is true
}

// Pivot groups ports from multiple snapshots by the chosen field.
func Pivot(snaps []Snapshot, opts PivotOptions) PivotResult {
	buckets := make(map[string][]scanner.Port)

	for _, s := range snaps {
		for _, p := range s.Ports {
			key := pivotKey(p, opts.Field)
			buckets[key] = append(buckets[key], p)
		}
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	return PivotResult{
		Field:   opts.Field,
		Buckets: buckets,
		Keys:    keys,
	}
}

func pivotKey(p scanner.Port, field PivotField) string {
	switch field {
	case PivotByProcess:
		if p.Process == "" {
			return "(unknown)"
		}
		return p.Process
	case PivotByPort:
		return itoa(p.Port)
	default: // PivotByProtocol
		if p.Protocol == "" {
			return "(unknown)"
		}
		return p.Protocol
	}
}
