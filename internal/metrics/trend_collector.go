package metrics

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// TrendCollector records gauge observations over time so that trend analysis
// can be performed on live metric data.
type TrendCollector struct {
	mu      sync.Mutex
	points  map[string][]snapshot.TrendPoint
	maxAge  time.Duration
}

// NewTrendCollector creates a TrendCollector that retains points up to maxAge.
// Pass 0 to keep all points.
func NewTrendCollector(maxAge time.Duration) *TrendCollector {
	return &TrendCollector{
		points: make(map[string][]snapshot.TrendPoint),
		maxAge: maxAge,
	}
}

// Record adds a new observation for the named metric at the current time.
func (tc *TrendCollector) Record(name string, value float64) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	now := time.Now().Unix()
	tc.points[name] = append(tc.points[name], snapshot.TrendPoint{
		Timestamp: now,
		Value:     value,
	})
	if tc.maxAge > 0 {
		cutoff := time.Now().Add(-tc.maxAge).Unix()
		filtered := tc.points[name][:0]
		for _, p := range tc.points[name] {
			if p.Timestamp >= cutoff {
				filtered = append(filtered, p)
			}
		}
		tc.points[name] = filtered
	}
}

// Points returns a copy of all recorded points for the given metric.
func (tc *TrendCollector) Points(name string) []snapshot.TrendPoint {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	src := tc.points[name]
	out := make([]snapshot.TrendPoint, len(src))
	copy(out, src)
	return out
}

// AllPoints returns copies of all metric point slices keyed by metric name.
func (tc *TrendCollector) AllPoints() map[string][]snapshot.TrendPoint {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	out := make(map[string][]snapshot.TrendPoint, len(tc.points))
	for k, v := range tc.points {
		cp := make([]snapshot.TrendPoint, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}
