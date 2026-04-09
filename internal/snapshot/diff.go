package snapshot

import (
	"github.com/user/portwatch/internal/scanner"
)

// DiffResult holds the outcome of comparing two port snapshots.
type DiffResult struct {
	Added   []scanner.Port
	Removed []scanner.Port
}

// Diff compares prev and curr port slices and returns added/removed ports.
func Diff(prev, curr []scanner.Port) DiffResult {
	prevIdx := indexPorts(prev)
	currIdx := indexPorts(curr)

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
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 10)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
