package snapshot

import (
	"fmt"
	"strconv"

	"github.com/user/portwatch/internal/scanner"
)

// DiffResult holds ports added or removed between two snapshots.
type DiffResult struct {
	Added   []scanner.PortState
	Removed []scanner.PortState
}

// Diff computes the added and removed ports between prev and current.
func Diff(prev, current []scanner.PortState) DiffResult {
	prevIdx := indexPorts(prev)
	currIdx := indexPorts(current)

	var result DiffResult
	for k, p := range currIdx {
		if _, ok := prevIdx[k]; !ok {
			result.Added = append(result.Added, p)
		}
	}
	for k, p := range prevIdx {
		if _, ok := currIdx[k]; !ok {
			result.Removed = append(result.Removed, p)
		}
	}
	return result
}

func indexPorts(ports []scanner.PortState) map[string]scanner.PortState {
	m := make(map[string]scanner.PortState, len(ports))
	for _, p := range ports {
		m[portKey(p)] = p
	}
	return m
}

func portKey(p scanner.PortState) string {
	return fmt.Sprintf("%s:%s", p.Protocol, itoa(p.Port))
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
