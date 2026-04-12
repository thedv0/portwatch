package snapshot

import (
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// GroupBy defines the field to group ports by.
type GroupBy string

const (
	GroupByProtocol GroupBy = "protocol"
	GroupByProcess  GroupBy = "process"
	GroupByPID      GroupBy = "pid"
	GroupByPort     GroupBy = "port"
)

// Group holds a named collection of ports.
type Group struct {
	Key   string
	Ports []scanner.Port
}

// GroupPorts partitions ports into named groups based on the given field.
// Groups are returned sorted by key for deterministic output.
func GroupPorts(ports []scanner.Port, by GroupBy) []Group {
	index := make(map[string][]scanner.Port)

	for _, p := range ports {
		var key string
		switch by {
		case GroupByProtocol:
			key = p.Protocol
			if key == "" {
				key = "unknown"
			}
		case GroupByProcess:
			key = p.Process
			if key == "" {
				key = "(unknown)"
			}
		case GroupByPID:
			key = itoa(p.PID)
		case GroupByPort:
			key = itoa(p.Port)
		default:
			key = "(all)"
		}
		index[key] = append(index[key], p)
	}

	groups := make([]Group, 0, len(index))
	for k, v := range index {
		groups = append(groups, Group{Key: k, Ports: v})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Key < groups[j].Key
	})
	return groups
}
