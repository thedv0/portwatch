package metrics

import "github.com/user/portwatch/internal/snapshot"

// DiffChainCollector records metrics derived from a diff chain.
type DiffChainCollector struct {
	reg          *Registry
	totalAdded   *Counter
	totalRemoved *Counter
	chainLength  *Gauge
}

// NewDiffChainCollector registers and returns a DiffChainCollector.
func NewDiffChainCollector(reg *Registry) *DiffChainCollector {
	return &DiffChainCollector{
		reg:          reg,
		totalAdded:   reg.Counter("diff_chain_added_total"),
		totalRemoved: reg.Counter("diff_chain_removed_total"),
		chainLength:  reg.Gauge("diff_chain_length"),
	}
}

// Collect updates metrics from the provided chain.
func (c *DiffChainCollector) Collect(chain []snapshot.ChainEntry) {
	c.chainLength.Set(float64(len(chain)))
	for _, e := range chain {
		c.totalAdded.Add(int64(len(e.Diff.Added)))
		c.totalRemoved.Add(int64(len(e.Diff.Removed)))
	}
}
