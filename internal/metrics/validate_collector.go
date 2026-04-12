package metrics

import (
	"github.com/user/portwatch/internal/snapshot"
)

// ValidateCollector records validation metrics into a Registry.
type ValidateCollector struct {
	reg          *Registry
	errorCount   *Counter
	warningCount *Counter
	infoCount    *Counter
	totalPorts   *Gauge
}

// NewValidateCollector creates a ValidateCollector backed by the given registry.
func NewValidateCollector(reg *Registry) *ValidateCollector {
	return &ValidateCollector{
		reg:          reg,
		errorCount:   reg.Counter("validate_errors_total"),
		warningCount: reg.Counter("validate_warnings_total"),
		infoCount:    reg.Counter("validate_info_total"),
		totalPorts:   reg.Gauge("validate_ports_total"),
	}
}

// Record updates counters and gauges from a ValidationResult and port count.
func (c *ValidateCollector) Record(result *snapshot.ValidationResult, portCount int) {
	c.totalPorts.Set(float64(portCount))
	for _, issue := range result.Issues {
		switch issue.Level {
		case snapshot.LevelError:
			c.errorCount.Inc()
		case snapshot.LevelWarning:
			c.warningCount.Inc()
		case snapshot.LevelInfo:
			c.infoCount.Inc()
		}
	}
}

// Reset clears all counters (useful between scan cycles).
func (c *ValidateCollector) Reset() {
	c.errorCount.Reset()
	c.warningCount.Reset()
	c.infoCount.Reset()
}
