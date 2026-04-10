package snapshot

import (
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// FilterOptions defines criteria for filtering a list of ports.
type FilterOptions struct {
	// Protocol filters by protocol ("tcp", "udp"); empty means all.
	Protocol string
	// MinPort is the lower bound (inclusive); 0 means no lower bound.
	MinPort uint16
	// MaxPort is the upper bound (inclusive); 0 means no upper bound.
	MaxPort uint16
	// PIDZeroOnly returns only entries where PID is 0 (unknown/kernel).
	PIDZeroOnly bool
	// ProcessName filters by process name substring (case-insensitive).
	ProcessName string
}

// Filter returns a new slice containing only the ports that match all
// non-zero criteria in opts.
func Filter(ports []scanner.PortState, opts FilterOptions) []scanner.PortState {
	result := make([]scanner.PortState, 0, len(ports))
	for _, p := range ports {
		if opts.Protocol != "" && !strings.EqualFold(p.Protocol, opts.Protocol) {
			continue
		}
		if opts.MinPort != 0 && p.Port < opts.MinPort {
			continue
		}
		if opts.MaxPort != 0 && p.Port > opts.MaxPort {
			continue
		}
		if opts.PIDZeroOnly && p.PID != 0 {
			continue
		}
		if opts.ProcessName != "" &&
			!strings.Contains(strings.ToLower(p.Process), strings.ToLower(opts.ProcessName)) {
			continue
		}
		result = append(result, p)
	}
	return result
}
