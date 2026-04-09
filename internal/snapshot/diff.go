package snapshot

import (
	"github.com/user/portwatch/internal/scanner"
)

// DiffResult holds the sets of added and removed ports between two snapshots.
type DiffResult struct {
	Added   []scanner.Port
	Removed []scanner.Port
}

// Diff compares prev and curr port slices and returns what was added or removed.
func Diff(prev, curr []scanner.Port) DiffResult {
	prevIdx := indexPorts(prev)
	currIdx := indexPorts(curr)

	var result DiffResult

	for k, p := range currIdx {
		if _, found := prevIdx[k]; !found {
			result.Added = append(result.Added, p)
		}
	}

	for k, p := range prevIdx {
		if _, found := currIdx[k]; !found {
			result.Removed = append(result.Removed, p)
		}
	}

	return result
}

func indexPorts(ports []scanner.Port) map[string]scanner.Port {
	idx := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		idx[portKey(p)] = p
	}
	return idx
}

func portKey(p scanner.Port) string {
	return p.Protocol + ":" + itoa(p.Port)
}

func itoa(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{digits[n%10]}, buf...)
		n /= 10
	}
	return string(buf)
}
