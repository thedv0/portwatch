package snapshot

import (
	"github.com/user/portwatch/internal/scanner"
)

// DiffResult holds the sets of ports added or removed between two snapshots.
type DiffResult struct {
	Added   []scanner.Port
	Removed []scanner.Port
}

// Diff computes the difference between a previous and current set of ports.
// Ports present in current but not prev are Added; ports in prev but not
// current are Removed.
func Diff(prev, current []scanner.Port) DiffResult {
	prevIdx := indexPorts(prev)
	currIdx := indexPorts(current)

	var result DiffResult

	for key, p := range currIdx {
		if _, exists := prevIdx[key]; !exists {
			result.Added = append(result.Added, p)
		}
	}

	for key, p := range prevIdx {
		if _, exists := currIdx[key]; !exists {
			result.Removed = append(result.Removed, p)
		}
	}

	return result
}

// indexPorts builds a map keyed by "proto:port" for fast lookup.
func indexPorts(ports []scanner.Port) map[string]scanner.Port {
	idx := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		idx[portKey(p)] = p
	}
	return idx
}

func portKey(p scanner.Port) string {
	return p.Protocol + ":" + itoa(p.Number)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 6)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
