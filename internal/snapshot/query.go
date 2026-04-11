package snapshot

import (
	"sort"
	"strings"
)

// SortField defines the field to sort ports by.
type SortField string

const (
	SortByPort     SortField = "port"
	SortByPID      SortField = "pid"
	SortByProtocol SortField = "protocol"
	SortByProcess  SortField = "process"
)

// QueryOptions controls sorting and pagination of port results.
type QueryOptions struct {
	SortBy    SortField
	Ascending bool
	Offset    int
	Limit     int
}

// DefaultQueryOptions returns sensible defaults.
func DefaultQueryOptions() QueryOptions {
	return QueryOptions{
		SortBy:    SortByPort,
		Ascending: true,
		Offset:    0,
		Limit:     0, // 0 means no limit
	}
}

// Query applies sorting and pagination to a slice of PortState.
func Query(ports []PortState, opts QueryOptions) []PortState {
	if len(ports) == 0 {
		return ports
	}

	result := make([]PortState, len(ports))
	copy(result, ports)

	sort.SliceStable(result, func(i, j int) bool {
		var less bool
		switch opts.SortBy {
		case SortByPID:
			less = result[i].PID < result[j].PID
		case SortByProtocol:
			less = strings.ToLower(result[i].Protocol) < strings.ToLower(result[j].Protocol)
		case SortByProcess:
			less = strings.ToLower(result[i].Process) < strings.ToLower(result[j].Process)
		default: // SortByPort
			less = result[i].Port < result[j].Port
		}
		if opts.Ascending {
			return less
		}
		return !less
	})

	if opts.Offset >= len(result) {
		return []PortState{}
	}
	result = result[opts.Offset:]

	if opts.Limit > 0 && opts.Limit < len(result) {
		result = result[:opts.Limit]
	}

	return result
}
