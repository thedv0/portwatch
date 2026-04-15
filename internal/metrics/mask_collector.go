package metrics

import (
	"github.com/user/portwatch/internal/snapshot"
)

// MaskCollector records metrics produced during a mask operation.
type MaskCollector struct {
	registry *Registry
	total    *Counter
	masked   *Counter
}

// NewMaskCollector creates a MaskCollector backed by the given Registry.
func NewMaskCollector(reg *Registry) *MaskCollector {
	return &MaskCollector{
		registry: reg,
		total:    reg.Counter("mask_total"),
		masked:   reg.Counter("mask_fields_masked"),
	}
}

// Collect records metrics for a completed mask operation. It counts the total
// number of entries processed and the number of fields that were masked based
// on the provided options.
func (c *MaskCollector) Collect(ports []snapshot.PortEntry, opts snapshot.MaskOptions) {
	n := int64(len(ports))
	c.total.Add(n)
	c.masked.Add(n * countMaskedFields(opts))
}

// countMaskedFields returns the number of fields that will be masked per entry
// given the provided MaskOptions.
func countMaskedFields(opts snapshot.MaskOptions) int64 {
	var count int64
	if opts.MaskProcess {
		count++
	}
	if opts.MaskPID {
		count++
	}
	if opts.MaskPort {
		count++
	}
	return count
}
