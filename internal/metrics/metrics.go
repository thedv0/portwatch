package metrics

import (
	"sync"
	"time"
)

// Counter tracks a named cumulative count.
type Counter struct {
	mu    sync.Mutex
	value int64
}

func (c *Counter) Inc() { c.mu.Lock(); c.value++; c.mu.Unlock() }
func (c *Counter) Add(n int64) { c.mu.Lock(); c.value += n; c.mu.Unlock() }
func (c *Counter) Value() int64 { c.mu.Lock(); defer c.mu.Unlock(); return c.value }
func (c *Counter) Reset() { c.mu.Lock(); c.value = 0; c.mu.Unlock() }

// Gauge holds a point-in-time value.
type Gauge struct {
	mu    sync.Mutex
	value float64
}

func (g *Gauge) Set(v float64) { g.mu.Lock(); g.value = v; g.mu.Unlock() }
func (g *Gauge) Value() float64 { g.mu.Lock(); defer g.mu.Unlock(); return g.value }

// Registry holds named metrics for the daemon.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
	gauges   map[string]*Gauge
	Started  time.Time
}

// New returns an initialised Registry.
func New() *Registry {
	return &Registry{
		counters: make(map[string]*Counter),
		gauges:   make(map[string]*Gauge),
		Started:  time.Now(),
	}
}

// Counter returns (creating if necessary) the named counter.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Gauge returns (creating if necessary) the named gauge.
func (r *Registry) Gauge(name string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := &Gauge{}
	r.gauges[name] = g
	return g
}

// Snapshot returns a copy of all current metric values.
func (r *Registry) Snapshot() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := map[string]interface{}{
		"uptime_seconds": time.Since(r.Started).Seconds(),
	}
	for k, c := range r.counters {
		out[k] = c.Value()
	}
	for k, g := range r.gauges {
		out[k] = g.Value()
	}
	return out
}
