package snapshot

import (
	"errors"
	"time"
)

// ThrottleOptions controls how port events are throttled to reduce noise.
type ThrottleOptions struct {
	// MinInterval is the minimum duration between repeated alerts for the same port key.
	MinInterval time.Duration
	// MaxBurst is the maximum number of events allowed before throttling kicks in.
	MaxBurst int
}

// DefaultThrottleOptions returns sensible defaults for throttling.
func DefaultThrottleOptions() ThrottleOptions {
	return ThrottleOptions{
		MinInterval: 30 * time.Second,
		MaxBurst:    3,
	}
}

// Validate returns an error if the options are invalid.
func (o ThrottleOptions) Validate() error {
	if o.MinInterval < 0 {
		return errors.New("throttle: MinInterval must be non-negative")
	}
	if o.MaxBurst < 1 {
		return errors.New("throttle: MaxBurst must be at least 1")
	}
	return nil
}

// throttleEntry tracks state for a single port key.
type throttleEntry struct {
	count    int
	lastSeen time.Time
}

// Throttler suppresses repeated port events within a configurable window.
type Throttler struct {
	opts    ThrottleOptions
	state   map[string]*throttleEntry
	clock   func() time.Time
}

// NewThrottler creates a Throttler with the given options.
func NewThrottler(opts ThrottleOptions) *Throttler {
	return &Throttler{
		opts:  opts,
		state: make(map[string]*throttleEntry),
		clock: time.Now,
	}
}

// Allow returns true if the event for the given key should be forwarded.
// It returns false when the key has been seen too recently or too frequently.
func (t *Throttler) Allow(key string) bool {
	now := t.clock()
	e, ok := t.state[key]
	if !ok {
		t.state[key] = &throttleEntry{count: 1, lastSeen: now}
		return true
	}
	if now.Sub(e.lastSeen) >= t.opts.MinInterval {
		e.count = 1
		e.lastSeen = now
		return true
	}
	if e.count < t.opts.MaxBurst {
		e.count++
		e.lastSeen = now
		return true
	}
	return false
}

// Reset clears throttle state for a specific key.
func (t *Throttler) Reset(key string) {
	delete(t.state, key)
}

// ResetAll clears all throttle state.
func (t *Throttler) ResetAll() {
	t.state = make(map[string]*throttleEntry)
}
