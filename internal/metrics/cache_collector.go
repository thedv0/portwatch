package metrics

import "sync/atomic"

// CacheCollector tracks cache hit/miss counters and integrates with the
// metrics Registry.
type CacheCollector struct {
	hits   atomic.Int64
	misses atomic.Int64
	reg    *Registry
}

// NewCacheCollector registers cache counters in reg and returns a collector.
func NewCacheCollector(reg *Registry) *CacheCollector {
	if reg != nil {
		reg.Counter("cache_hits_total")
		reg.Counter("cache_misses_total")
	}
	return &CacheCollector{reg: reg}
}

// RecordHit increments the hit counter.
func (c *CacheCollector) RecordHit() {
	c.hits.Add(1)
	if c.reg != nil {
		c.reg.Counter("cache_hits_total").Inc()
	}
}

// RecordMiss increments the miss counter.
func (c *CacheCollector) RecordMiss() {
	c.misses.Add(1)
	if c.reg != nil {
		c.reg.Counter("cache_misses_total").Inc()
	}
}

// Hits returns the total number of cache hits recorded.
func (c *CacheCollector) Hits() int64 {
	return c.hits.Load()
}

// Misses returns the total number of cache misses recorded.
func (c *CacheCollector) Misses() int64 {
	return c.misses.Load()
}

// Reset zeroes both counters (useful in tests).
func (c *CacheCollector) Reset() {
	c.hits.Store(0)
	c.misses.Store(0)
}
