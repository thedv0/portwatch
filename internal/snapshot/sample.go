package snapshot

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// SampleOptions controls how port samples are collected and stored.
type SampleOptions struct {
	// MaxSamples is the maximum number of samples to retain in a window.
	MaxSamples int
	// Interval is the minimum duration between samples.
	Interval time.Duration
}

// DefaultSampleOptions returns sensible defaults for sampling.
func DefaultSampleOptions() SampleOptions {
	return SampleOptions{
		MaxSamples: 60,
		Interval:   30 * time.Second,
	}
}

// Sample holds a single point-in-time capture of open ports.
type Sample struct {
	CapturedAt time.Time
	Ports      []scanner.Port
}

// SampleWindow holds an ordered, bounded collection of samples.
type SampleWindow struct {
	opts    SampleOptions
	samples []Sample
}

// NewSampleWindow creates a SampleWindow with the given options.
func NewSampleWindow(opts SampleOptions) *SampleWindow {
	if opts.MaxSamples <= 0 {
		opts.MaxSamples = DefaultSampleOptions().MaxSamples
	}
	if opts.Interval <= 0 {
		opts.Interval = DefaultSampleOptions().Interval
	}
	return &SampleWindow{opts: opts}
}

// Add appends a new sample to the window, enforcing interval and max size.
// Returns true if the sample was accepted, false if it was too soon.
func (w *SampleWindow) Add(ports []scanner.Port, at time.Time) bool {
	if len(w.samples) > 0 {
		last := w.samples[len(w.samples)-1]
		if at.Sub(last.CapturedAt) < w.opts.Interval {
			return false
		}
	}
	w.samples = append(w.samples, Sample{CapturedAt: at, Ports: ports})
	if len(w.samples) > w.opts.MaxSamples {
		w.samples = w.samples[len(w.samples)-w.opts.MaxSamples:]
	}
	return true
}

// All returns a copy of all samples in the window, oldest first.
func (w *SampleWindow) All() []Sample {
	out := make([]Sample, len(w.samples))
	copy(out, w.samples)
	return out
}

// Len returns the number of samples currently held.
func (w *SampleWindow) Len() int { return len(w.samples) }

// Latest returns the most recent sample, or false if the window is empty.
func (w *SampleWindow) Latest() (Sample, bool) {
	if len(w.samples) == 0 {
		return Sample{}, false
	}
	return w.samples[len(w.samples)-1], true
}
