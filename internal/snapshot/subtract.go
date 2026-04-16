package snapshot

import (
	"github.com/user/portwatch/internal/scanner"
)

// SubtractOptions configures the Subtract operation.
type SubtractOptions struct {
	// ByPort removes ports found in the right set by port number only.
	ByPort bool
	// ByProcess removes ports found in the right set by process name only.
	ByProcess bool
}

// DefaultSubtractOptions returns sensible defaults: subtract by port+protocol key.
func DefaultSubtractOptions() SubtractOptions {
	return SubtractOptions{}
}

// subtractKey builds a lookup key from a port entry.
func subtractKey(p scanner.Port, opts SubtractOptions) string {
	switch {
	case opts.ByPort:
		return itoa(p.Port)
	case opts.ByProcess:
		return p.Process
	default:
		return p.Protocol + ":" + itoa(p.Port)
	}
}

// Subtract returns ports in left that are NOT present in right, according to opts.
// The right slice acts as the exclusion set.
func Subtract(left, right []scanner.Port, opts SubtractOptions) []scanner.Port {
	if len(right) == 0 {
		return left
	}

	exclude := make(map[string]struct{}, len(right))
	for _, p := range right {
		exclude[subtractKey(p, opts)] = struct{}{}
	}

	result := make([]scanner.Port, 0, len(left))
	for _, p := range left {
		if _, found := exclude[subtractKey(p, opts)]; !found {
			result = append(result, p)
		}
	}
	return result
}
