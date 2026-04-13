package metrics

import "github.com/user/portwatch/internal/snapshot"

// ThrottleCollector records throttle-related metrics into a Registry.
type ThrottleCollector struct {
	allowed *Counter
	blocked *Counter
	active  *Gauge
}

// NewThrottleCollector registers and returns a ThrottleCollector.
func NewThrottleCollector(reg *Registry) *ThrottleCollector {
	return &ThrottleCollector{
		allowed: reg.Counter("throttle_allowed_total"),
		blocked: reg.Counter("throttle_blocked_total"),
		active:  reg.Gauge("throttle_active_keys"),
	}
}

// RecordAllow increments the allowed counter.
func (c *ThrottleCollector) RecordAllow() {
	c.allowed.Inc()
}

// RecordBlock increments the blocked counter.
func (c *ThrottleCollector) RecordBlock() {
	c.blocked.Inc()
}

// SetActiveKeys updates the gauge tracking how many keys are being tracked.
func (c *ThrottleCollector) SetActiveKeys(n int) {
	c.active.Set(float64(n))
}

// Observe wraps a Throttler.Allow call and records the outcome automatically.
func (c *ThrottleCollector) Observe(th *snapshot.Throttler, key string) bool {
	ok := th.Allow(key)
	if ok {
		c.allowed.Inc()
	} else {
		c.blocked.Inc()
	}
	return ok
}
