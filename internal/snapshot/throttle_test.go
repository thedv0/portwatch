package snapshot

import (
	"testing"
	"time"
)

func TestDefaultThrottleOptions_Values(t *testing.T) {
	opts := DefaultThrottleOptions()
	if opts.MinInterval <= 0 {
		t.Error("expected positive MinInterval")
	}
	if opts.MaxBurst < 1 {
		t.Error("expected MaxBurst >= 1")
	}
}

func TestValidate_ThrottleOptions_Valid(t *testing.T) {
	if err := DefaultThrottleOptions().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_ThrottleOptions_NegativeInterval(t *testing.T) {
	opts := DefaultThrottleOptions()
	opts.MinInterval = -1
	if err := opts.Validate(); err == nil {
		t.Error("expected error for negative MinInterval")
	}
}

func TestValidate_ThrottleOptions_ZeroBurst(t *testing.T) {
	opts := DefaultThrottleOptions()
	opts.MaxBurst = 0
	if err := opts.Validate(); err == nil {
		t.Error("expected error for zero MaxBurst")
	}
}

func TestThrottler_FirstEventAllowed(t *testing.T) {
	th := NewThrottler(DefaultThrottleOptions())
	if !th.Allow("tcp:8080") {
		t.Error("first event should be allowed")
	}
}

func TestThrottler_BurstAllowed(t *testing.T) {
	opts := ThrottleOptions{MinInterval: 10 * time.Second, MaxBurst: 3}
	th := NewThrottler(opts)
	th.clock = func() time.Time { return time.Unix(1000, 0) }
	for i := 0; i < 3; i++ {
		if !th.Allow("tcp:9090") {
			t.Fatalf("event %d should be allowed within burst", i+1)
		}
	}
	if th.Allow("tcp:9090") {
		t.Error("event beyond burst should be throttled")
	}
}

func TestThrottler_ResetsAfterInterval(t *testing.T) {
	opts := ThrottleOptions{MinInterval: 5 * time.Second, MaxBurst: 1}
	now := time.Unix(1000, 0)
	th := NewThrottler(opts)
	th.clock = func() time.Time { return now }
	th.Allow("tcp:443")
	if th.Allow("tcp:443") {
		t.Error("second event within interval should be throttled")
	}
	now = now.Add(6 * time.Second)
	if !th.Allow("tcp:443") {
		t.Error("event after interval should be allowed")
	}
}

func TestThrottler_Reset_ClearsKey(t *testing.T) {
	opts := ThrottleOptions{MinInterval: 60 * time.Second, MaxBurst: 1}
	th := NewThrottler(opts)
	th.Allow("tcp:22")
	th.Reset("tcp:22")
	if !th.Allow("tcp:22") {
		t.Error("event after Reset should be allowed")
	}
}

func TestThrottler_ResetAll(t *testing.T) {
	opts := ThrottleOptions{MinInterval: 60 * time.Second, MaxBurst: 1}
	th := NewThrottler(opts)
	th.Allow("tcp:22")
	th.Allow("udp:53")
	th.ResetAll()
	if !th.Allow("tcp:22") || !th.Allow("udp:53") {
		t.Error("events after ResetAll should be allowed")
	}
}
