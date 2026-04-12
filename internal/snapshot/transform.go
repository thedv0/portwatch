package snapshot

import (
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// TransformOptions controls how ports are transformed.
type TransformOptions struct {
	// MapProtocol replaces protocol values using the provided map.
	MapProtocol map[string]string
	// RenameProcess applies a rename map to process names (exact match).
	RenameProcess map[string]string
	// OffsetPorts adds a fixed integer offset to every port number (useful for testing).
	OffsetPorts int
	// ForceUpperProtocol uppercases all protocol strings.
	ForceUpperProtocol bool
}

// DefaultTransformOptions returns a no-op TransformOptions.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{}
}

// Transform applies the given options to a slice of ports, returning a new
// slice with all transformations applied. The original slice is not modified.
func Transform(ports []scanner.Port, opts TransformOptions) []scanner.Port {
	out := make([]scanner.Port, len(ports))
	for i, p := range ports {
		cloned := p

		// Protocol mapping takes precedence over ForceUpperProtocol.
		if opts.MapProtocol != nil {
			if mapped, ok := opts.MapProtocol[strings.ToLower(cloned.Protocol)]; ok {
				cloned.Protocol = mapped
			}
		}
		if opts.ForceUpperProtocol {
			cloned.Protocol = strings.ToUpper(cloned.Protocol)
		}

		// Process rename.
		if opts.RenameProcess != nil {
			if renamed, ok := opts.RenameProcess[cloned.Process]; ok {
				cloned.Process = renamed
			}
		}

		// Port offset — clamp to valid range [1, 65535].
		if opts.OffsetPorts != 0 {
			newPort := cloned.Port + opts.OffsetPorts
			if newPort < 1 {
				newPort = 1
			}
			if newPort > 65535 {
				newPort = 65535
			}
			cloned.Port = newPort
		}

		out[i] = cloned
	}
	return out
}
