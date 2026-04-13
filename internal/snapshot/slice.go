package snapshot

import "github.com/wricardo/portwatch/internal/scanner"

// SliceOptions controls how a port list is sliced.
type SliceOptions struct {
	// Offset is the number of entries to skip from the start.
	Offset int
	// Limit is the maximum number of entries to return. 0 means no limit.
	Limit int
	// Reverse reverses the order before slicing.
	Reverse bool
}

// DefaultSliceOptions returns a SliceOptions with no-op defaults.
func DefaultSliceOptions() SliceOptions {
	return SliceOptions{
		Offset:  0,
		Limit:   0,
		Reverse: false,
	}
}

// Slice applies offset, limit, and optional reversal to a port list.
func Slice(ports []scanner.Port, opts SliceOptions) []scanner.Port {
	if len(ports) == 0 {
		return []scanner.Port{}
	}

	result := make([]scanner.Port, len(ports))
	copy(result, ports)

	if opts.Reverse {
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
	}

	if opts.Offset >= len(result) {
		return []scanner.Port{}
	}
	result = result[opts.Offset:]

	if opts.Limit > 0 && opts.Limit < len(result) {
		result = result[:opts.Limit]
	}

	return result
}
