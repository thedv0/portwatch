package snapshot

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSamplePorts(nums ...int) []scanner.Port {
	ports := make([]scanner.Port, len(nums))
	for i, n := range nums {
		ports[i] = scanner.Port{Port: n, Protocol: "tcp"}
	}
	return ports
}

func TestDefaultSampleOptions_Values(t *testing.T) {
	opts := DefaultSampleOptions()
	if opts.MaxSamples != 60 {
		t.Errorf("expected MaxSamples 60, got %d", opts.MaxSamples)
	}
	if opts.Interval != 30*time.Second {
		t.Errorf("expected Interval 30s, got %v", opts.Interval)
	}
}

func TestSampleWindow_Add_AcceptsSample(t *testing.T) {
	w := NewSampleWindow(SampleOptions{MaxSamples: 10, Interval: time.Second})
	base := time.Now()
	ok := w.Add(makeSamplePorts(80), base)
	if !ok {
		t.Fatal("expected first sample to be accepted")
	}
	if w.Len() != 1 {
		t.Fatalf("expected len 1, got %d", w.Len())
	}
}

func TestSampleWindow_Add_RejectsTooSoon(t *testing.T) {
	w := NewSampleWindow(SampleOptions{MaxSamples: 10, Interval: 5 * time.Second})
	base := time.Now()
	w.Add(makeSamplePorts(80), base)
	ok := w.Add(makeSamplePorts(443), base.Add(time.Second))
	if ok {
		t.Fatal("expected sample within interval to be rejected")
	}
	if w.Len() != 1 {
		t.Fatalf("expected len 1, got %d", w.Len())
	}
}

func TestSampleWindow_Add_EnforcesMaxSamples(t *testing.T) {
	w := NewSampleWindow(SampleOptions{MaxSamples: 3, Interval: time.Millisecond})
	base := time.Now()
	for i := 0; i < 5; i++ {
		w.Add(makeSamplePorts(i+1), base.Add(time.Duration(i)*time.Second))
	}
	if w.Len() != 3 {
		t.Fatalf("expected len 3, got %d", w.Len())
	}
}

func TestSampleWindow_Latest_Empty(t *testing.T) {
	w := NewSampleWindow(DefaultSampleOptions())
	_, ok := w.Latest()
	if ok {
		t.Fatal("expected Latest to return false on empty window")
	}
}

func TestSampleWindow_Latest_ReturnsMostRecent(t *testing.T) {
	w := NewSampleWindow(SampleOptions{MaxSamples: 10, Interval: time.Millisecond})
	base := time.Now()
	w.Add(makeSamplePorts(80), base)
	w.Add(makeSamplePorts(443), base.Add(time.Second))
	s, ok := w.Latest()
	if !ok {
		t.Fatal("expected Latest to return a sample")
	}
	if len(s.Ports) != 1 || s.Ports[0].Port != 443 {
		t.Errorf("expected latest port 443, got %+v", s.Ports)
	}
}

func TestSampleWindow_All_ReturnsCopy(t *testing.T) {
	w := NewSampleWindow(SampleOptions{MaxSamples: 10, Interval: time.Millisecond})
	base := time.Now()
	w.Add(makeSamplePorts(80), base)
	w.Add(makeSamplePorts(443), base.Add(time.Second))
	all := w.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 samples, got %d", len(all))
	}
	// Mutating the copy must not affect the window.
	all[0].Ports = nil
	if w.Len() != 2 {
		t.Error("mutating All() result should not affect window")
	}
}
