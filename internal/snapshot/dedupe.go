package snapshot

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// DedupeOptions controls how deduplication is performed.
type DedupeOptions struct {
	// PreferHigherPID keeps the entry with the higher PID when duplicates exist.
	PreferHigherPID bool
	// IgnoreProcess ignores the process name when determining duplicates.
	IgnoreProcess bool
}

// DefaultDedupeOptions returns sensible defaults for deduplication.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		PreferHigherPID: false,
		IgnoreProcess:   false,
	}
}

// dedupeKey builds a string key for a port entry based on options.
func dedupeKey(p scanner.Port, opts DedupeOptions) string {
	if opts.IgnoreProcess {
		return fmt.Sprintf("%s:%d", p.Protocol, p.Port)
	}
	return fmt.Sprintf("%s:%d:%s", p.Protocol, p.Port, p.Process)
}

// Dedupe removes duplicate port entries from the slice.
// When duplicates are found, the first occurrence is kept unless
// PreferHigherPID is set, in which case the entry with the higher PID wins.
func Dedupe(ports []scanner.Port, opts DedupeOptions) []scanner.Port {
	if len(ports) == 0 {
		return ports
	}

	seen := make(map[string]int) // key -> index in result
	result := make([]scanner.Port, 0, len(ports))

	for _, p := range ports {
		key := dedupeKey(p, opts)
		if idx, exists := seen[key]; exists {
			if opts.PreferHigherPID && p.PID > result[idx].PID {
				result[idx] = p
			}
			continue
		}
		seen[key] = len(result)
		result = append(result, p)
	}

	return result
}
