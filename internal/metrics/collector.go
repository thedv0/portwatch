package metrics

import "github.com/user/portwatch/internal/scanner"

// PortStats groups the values collected during a single scan cycle.
type PortStats struct {
	OpenPorts int
	Alerts    int
	ScanError bool
}

// Collector wraps a Registry and provides domain-specific recording helpers.
type Collector struct {
	reg        *Registry
	ScanTotal  *Counter
	ScanErrors *Counter
	AlertsSent *Counter
	OpenPorts  *Gauge
}

// NewCollector creates a Collector backed by reg.
func NewCollector(reg *Registry) *Collector {
	return &Collector{
		reg:        reg,
		ScanTotal:  reg.Counter("scans_total"),
		ScanErrors: reg.Counter("scan_errors_total"),
		AlertsSent: reg.Counter("alerts_sent_total"),
		OpenPorts:  reg.Gauge("open_ports"),
	}
}

// Record updates all metrics from a completed scan cycle.
func (c *Collector) Record(stats PortStats) {
	c.ScanTotal.Inc()
	if stats.ScanError {
		c.ScanErrors.Inc()
	}
	c.AlertsSent.Add(int64(stats.Alerts))
	c.OpenPorts.Set(float64(stats.OpenPorts))
}

// RecordPorts is a convenience that counts open ports from a scanner result.
func (c *Collector) RecordPorts(ports []scanner.Port, alerts int, scanErr bool) {
	c.Record(PortStats{
		OpenPorts: len(ports),
		Alerts:    alerts,
		ScanError: scanErr,
	})
}
