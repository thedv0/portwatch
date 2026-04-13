package metrics

import "github.com/yourorg/portwatch/internal/snapshot"

// AnomalyCollector records anomaly detection metrics into a Registry.
type AnomalyCollector struct {
	reg        *Registry
	total      *Counter
	newPorts   *Counter
	gonePorts  *Counter
	pidChanges *Counter
	spikes     *Counter
}

// NewAnomalyCollector creates an AnomalyCollector backed by reg.
func NewAnomalyCollector(reg *Registry) *AnomalyCollector {
	return &AnomalyCollector{
		reg:        reg,
		total:      reg.Counter("anomaly_total"),
		newPorts:   reg.Counter("anomaly_new_ports"),
		gonePorts:  reg.Counter("anomaly_gone_ports"),
		pidChanges: reg.Counter("anomaly_pid_changes"),
		spikes:     reg.Counter("anomaly_spikes"),
	}
}

// Record updates counters based on the provided anomalies.
func (c *AnomalyCollector) Record(anomalies []snapshot.Anomaly) {
	for _, a := range anomalies {
		c.total.Inc()
		switch a.Type {
		case snapshot.AnomalyNewPort:
			c.newPorts.Inc()
		case snapshot.AnomalyGonePort:
			c.gonePorts.Inc()
		case snapshot.AnomalyPIDChanged:
			c.pidChanges.Inc()
		case snapshot.AnomalyPortSpike:
			c.spikes.Inc()
		}
	}
}
