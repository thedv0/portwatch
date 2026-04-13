package metrics

import "github.com/user/portwatch/internal/snapshot"

// ReplayCollector records metrics derived from a replay sequence.
type ReplayCollector struct {
	framesTotal  *Counter
	portsAdded   *Counter
	portsRemoved *Counter
	maxOpen      *Gauge
}

// NewReplayCollector registers replay metrics on the given registry.
func NewReplayCollector(reg *Registry) *ReplayCollector {
	return &ReplayCollector{
		framesTotal:  reg.Counter("replay_frames_total"),
		portsAdded:   reg.Counter("replay_ports_added_total"),
		portsRemoved: reg.Counter("replay_ports_removed_total"),
		maxOpen:      reg.Gauge("replay_max_open_ports"),
	}
}

// Collect updates counters and gauges from the provided replay frames.
func (c *ReplayCollector) Collect(frames []snapshot.ReplayFrame) {
	var maxOpen int
	for _, f := range frames {
		c.framesTotal.Inc()
		c.portsAdded.Add(int64(len(f.Diff.Added)))
		c.portsRemoved.Add(int64(len(f.Diff.Removed)))
		if len(f.Ports) > maxOpen {
			maxOpen = len(f.Ports)
		}
	}
	if len(frames) > 0 {
		c.maxOpen.Set(float64(maxOpen))
	}
}
