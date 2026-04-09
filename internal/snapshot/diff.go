package snapshot

import "github.com/user/portwatch/internal/scanner"

// ChangeKind describes whether a port appeared or disappeared.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
)

// Change records a single port that appeared or disappeared between snapshots.
type Change struct {
	Kind ChangeKind   `json:"kind"`
	Port scanner.Port `json:"port"`
}

// Diff computes the ports that were added or removed between two port lists.
// prev is the older list; curr is the current list.
func Diff(prev, curr []scanner.Port) []Change {
	prevSet := indexPorts(prev)
	currSet := indexPorts(curr)

	var changes []Change

	for key, p := range currSet {
		if _, exists := prevSet[key]; !exists {
			changes = append(changes, Change{Kind: Added, Port: p})
		}
	}

	for key, p := range prevSet {
		if _, exists := currSet[key]; !exists {
			changes = append(changes, Change{Kind: Removed, Port: p})
		}
	}

	return changes
}

func indexPorts(ports []scanner.Port) map[string]scanner.Port {
	m := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		key := portKey(p)
		m[key] = p
	}
	return m
}

func portKey(p scanner.Port) string {
	return p.Protocol + ":" + itoa(p.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
