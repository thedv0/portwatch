package snapshot

import (
	"errors"
	"sort"
	"time"
)

// EvictPolicy controls how ports are evicted from a working set.
type EvictPolicy int

const (
	EvictByAge      EvictPolicy = iota // remove ports not seen since a cutoff
	EvictByCount                       // keep only the N most-recently-seen ports
	EvictByIdleTime                    // remove ports idle longer than a duration
)

// DefaultEvictOptions returns sensible defaults.
func DefaultEvictOptions() EvictOptions {
	return EvictOptions{
		Policy:   EvictByAge,
		MaxAge:   30 * time.Minute,
		MaxCount: 1000,
		IdleTime: 10 * time.Minute,
	}
}

// EvictOptions configures the Evict function.
type EvictOptions struct {
	Policy   EvictPolicy
	MaxAge   time.Duration // used by EvictByAge
	MaxCount int           // used by EvictByCount
	IdleTime time.Duration // used by EvictByIdleTime
	Now      func() time.Time
}

func (o *EvictOptions) now() time.Time {
	if o.Now != nil {
		return o.Now()
	}
	return time.Now()
}

// Evict removes ports from ports according to the configured policy.
// lastSeen maps portKey(p) -> last observed timestamp.
func Evict(ports []PortState, lastSeen map[string]time.Time, opts EvictOptions) ([]PortState, error) {
	switch opts.Policy {
	case EvictByAge:
		return evictByAge(ports, lastSeen, opts)
	case EvictByCount:
		return evictByCount(ports, lastSeen, opts)
	case EvictByIdleTime:
		return evictByIdle(ports, lastSeen, opts)
	default:
		return nil, errors.New("evict: unknown policy")
	}
}

func evictByAge(ports []PortState, lastSeen map[string]time.Time, opts EvictOptions) ([]PortState, error) {
	cutoff := opts.now().Add(-opts.MaxAge)
	out := ports[:0:len(ports)]
	for _, p := range ports {
		if t, ok := lastSeen[evictKey(p)]; !ok || !t.Before(cutoff) {
			out = append(out, p)
		}
	}
	return out, nil
}

func evictByCount(ports []PortState, lastSeen map[string]time.Time, opts EvictOptions) ([]PortState, error) {
	if opts.MaxCount <= 0 {
		return nil, errors.New("evict: MaxCount must be positive")
	}
	if len(ports) <= opts.MaxCount {
		return ports, nil
	}
	// sort oldest-first, keep the newest MaxCount
	sorted := make([]PortState, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		ti := lastSeen[evictKey(sorted[i])]
		tj := lastSeen[evictKey(sorted[j])]
		return ti.Before(tj)
	})
	return sorted[len(sorted)-opts.MaxCount:], nil
}

func evictByIdle(ports []PortState, lastSeen map[string]time.Time, opts EvictOptions) ([]PortState, error) {
	cutoff := opts.now().Add(-opts.IdleTime)
	out := ports[:0:len(ports)]
	for _, p := range ports {
		if t, ok := lastSeen[evictKey(p)]; !ok || !t.Before(cutoff) {
			out = append(out, p)
		}
	}
	return out, nil
}

func evictKey(p PortState) string {
	return p.Protocol + ":" + itoa(p.Port)
}
