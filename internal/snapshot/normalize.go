package snapshot

import (
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// NormalizeOptions controls how port normalization is applied.
type NormalizeOptions struct {
	// LowercaseProtocol converts protocol strings to lowercase.
	LowercaseProtocol bool
	// TrimProcessName trims whitespace from process names.
	TrimProcessName bool
	// ZeroInvalidPID sets negative PIDs to zero.
	ZeroInvalidPID bool
	// ClampPort clamps port numbers to the valid range [0, 65535].
	ClampPort bool
}

// DefaultNormalizeOptions returns sensible defaults for normalization.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		LowercaseProtocol: true,
		TrimProcessName:   true,
		ZeroInvalidPID:    true,
		ClampPort:         true,
	}
}

// Normalize applies normalization rules to a slice of ports and returns
// a new slice with the corrections applied. The original slice is not
// modified.
func Normalize(ports []scanner.Port, opts NormalizeOptions) []scanner.Port {
	out := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		np := p
		if opts.LowercaseProtocol {
			np.Protocol = strings.ToLower(np.Protocol)
		}
		if opts.TrimProcessName {
			np.Process = strings.TrimSpace(np.Process)
		}
		if opts.ZeroInvalidPID && np.PID < 0 {
			np.PID = 0
		}
		if opts.ClampPort {
			if np.Port < 0 {
				np.Port = 0
			} else if np.Port > 65535 {
				np.Port = 65535
			}
		}
		out = append(out, np)
	}
	return out
}
