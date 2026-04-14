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

	var fieldsPerEntry int64
	if opts.MaskProcess {
		fieldsPerEntry++
	}
	if opts.MaskPID {
		fieldsPerEntry++
	}
	if opts.MaskPort {
		fieldsPerEntry++
	}
	c.masked.Add(n * fieldsPerEntry)
}
