package metrics

import (
	"github.com/user/portwatch/internal/report"
)

// NewUnionCollector records metrics from a UnionReport into the given Registry.
// It tracks the number of input snapshots, total ports after union, and whether
// deduplication was enabled.
func NewUnionCollector(reg *Registry) *UnionCollector {
	return &UnionCollector{
		inputSnaps: reg.Counter("union_input_snaps_total"),
		totalPorts: reg.Gauge("union_total_ports"),
		dedupEnabled: reg.Gauge("union_dedup_enabled"),
	}
}

// UnionCollector holds the metric handles for union operations.
type UnionCollector struct {
	inputSnaps   Counter
	totalPorts   Gauge
	dedupEnabled Gauge
}

// Collect updates all metrics from the given UnionReport.
func (c *UnionCollector) Collect(r report.UnionReport) {
	c.inputSnaps.Add(int64(r.InputSnaps))
	c.totalPorts.Set(float64(r.TotalPorts))
	if r.DedupEnabled {
		c.dedupEnabled.Set(1)
	} else {
		c.dedupEnabled.Set(0)
	}
}
