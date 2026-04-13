package metrics

import (
	"github.com/netwatch/portwatch/internal/scanner"
)

// ClusterCollector records metrics derived from a cluster result.
type ClusterCollector struct {
	reg          *Registry
	clusterCount *Counter
	totalPorts   *Counter
	largestSize  *Gauge
}

// NewClusterCollector registers and returns a ClusterCollector.
func NewClusterCollector(reg *Registry) *ClusterCollector {
	return &ClusterCollector{
		reg:          reg,
		clusterCount: reg.Counter("cluster_count"),
		totalPorts:   reg.Counter("cluster_total_ports"),
		largestSize:  reg.Gauge("cluster_largest_size"),
	}
}

// Observe records metrics for the given cluster result.
func (c *ClusterCollector) Observe(clusters map[string][]scanner.Port) {
	c.clusterCount.Reset()
	c.totalPorts.Reset()

	var largest int
	for _, ports := range clusters {
		c.clusterCount.Inc()
		for range ports {
			c.totalPorts.Inc()
		}
		if len(ports) > largest {
			largest = len(ports)
		}
	}
	c.largestSize.Set(float64(largest))
}
