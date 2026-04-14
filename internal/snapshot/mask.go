package snapshot

import (
	"strings"
)

// MaskOptions controls which fields are masked in port entries.
type MaskOptions struct {
	// MaskProcess replaces process names with a redacted placeholder.
	MaskProcess bool
	// MaskPID zeroes out PID values.
	MaskPID bool
	// MaskPort replaces port numbers with zero.
	MaskPort bool
	// Placeholder is the string used for masked string fields.
	Placeholder string
}

// DefaultMaskOptions returns sensible defaults: only process names are masked.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		MaskProcess: true,
		MaskPID:     false,
		MaskPort:    false,
		Placeholder: "[redacted]",
	}
}

// Mask applies the given MaskOptions to a slice of PortEntry, returning a
// new slice with sensitive fields replaced. The original slice is not modified.
func Mask(ports []PortEntry, opts MaskOptions) []PortEntry {
	ph := opts.Placeholder
	if strings.TrimSpace(ph) == "" {
		ph = "[redacted]"
	}

	out := make([]PortEntry, len(ports))
	for i, p := range ports {
		cp := p
		if opts.MaskProcess {
			cp.Process = ph
		}
		if opts.MaskPID {
			cp.PID = 0
		}
		if opts.MaskPort {
			cp.Port = 0
		}
		out[i] = cp
	}
	return out
}
