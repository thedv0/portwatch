package snapshot

import (
	"github.com/user/portwatch/internal/scanner"
)

// DiffResult holds the ports added and removed between two snapshots.
type DiffResult struct {
	Added   []scanner.Port
	Removed []scanner.Port
}

// Diff compares prev and curr port slices and returns what changed.
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
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
