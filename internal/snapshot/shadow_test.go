package snapshot

import (
	"testing"
	"time"
)

func shPort(proto string, port int, process string) PortState {
	return PortState{Protocol: proto, Port: port, Process: process}
}

var fixedShadowClock = func() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestShadow_NoDivergence(t *testing.T) {
	primary := []PortState{shPort("tcp", 80, "nginx"), shPort("tcp", 443, "nginx")}
	shadow := []PortState{shPort("tcp", 80, "nginx"), shPort("tcp", 443, "nginx")}
	opts := DefaultShadowOptions()
	opts.Clock = fixedShadowClock
	res, err := Shadow(primary, shadow, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Diverged {
		t.Errorf("expected no divergence")
	}
	if len(res.Divergences) != 0 {
		t.Errorf("expected 0 divergences, got %d", len(res.Divergences))
	}
}

func TestShadow_MissingInPrimary(t *testing.T) {
	primary := []PortState{shPort("tcp", 80, "nginx")}
	shadow := []PortState{shPort("tcp", 80, "nginx"), shPort("tcp", 8080, "proxy")}
	opts := DefaultShadowOptions()
	opts.Clock = fixedShadowClock
	res, err := Shadow(primary, shadow, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Diverged {
		t.Errorf("expected divergence")
	}
	if len(res.Divergences) != 1 {
		t.Errorf("expected 1 divergence, got %d", len(res.Divergences))
	}
}

func TestShadow_MissingInShadow(t *testing.T) {
	primary := []PortState{shPort("tcp", 80, "nginx"), shPort("udp", 53, "dns")}
	shadow := []PortState{shPort("tcp", 80, "nginx")}
	opts := DefaultShadowOptions()
	opts.Clock = fixedShadowClock
	res, _ := Shadow(primary, shadow, opts)
	if !res.Diverged {
		t.Errorf("expected divergence")
	}
}

func TestShadow_ProcessMismatch(t *testing.T) {
	primary := []PortState{shPort("tcp", 80, "apache")}
	shadow := []PortState{shPort("tcp", 80, "nginx")}
	opts := DefaultShadowOptions()
	opts.Clock = fixedShadowClock
	res, _ := Shadow(primary, shadow, opts)
	if !res.Diverged {
		t.Errorf("expected divergence on process mismatch")
	}
}

func TestShadow_IgnoreProcess(t *testing.T) {
	primary := []PortState{shPort("tcp", 80, "apache")}
	shadow := []PortState{shPort("tcp", 80, "nginx")}
	opts := DefaultShadowOptions()
	opts.Clock = fixedShadowClock
	opts.IgnoreProcess = true
	res, _ := Shadow(primary, shadow, opts)
	if res.Diverged {
		t.Errorf("expected no divergence when IgnoreProcess=true")
	}
}

func TestShadow_Tolerance(t *testing.T) {
	primary := []PortState{shPort("tcp", 80, "nginx")}
	shadow := []PortState{shPort("tcp", 80, "nginx"), shPort("tcp", 9090, "extra")}
	opts := DefaultShadowOptions()
	opts.Clock = fixedShadowClock
	opts.Tolerance = 1
	res, _ := Shadow(primary, shadow, opts)
	if res.Diverged {
		t.Errorf("expected no divergence within tolerance")
	}
}

func TestShadow_NilClock_ReturnsError(t *testing.T) {
	opts := DefaultShadowOptions()
	opts.Clock = nil
	_, err := Shadow(nil, nil, opts)
	if err == nil {
		t.Errorf("expected error for nil Clock")
	}
}

func TestShadow_LengthsRecorded(t *testing.T) {
	primary := []PortState{shPort("tcp", 80, "nginx"), shPort("tcp", 443, "nginx")}
	shadow := []PortState{shPort("tcp", 80, "nginx")}
	opts := DefaultShadowOptions()
	opts.Clock = fixedShadowClock
	res, _ := Shadow(primary, shadow, opts)
	if res.PrimaryLen != 2 {
		t.Errorf("expected PrimaryLen=2, got %d", res.PrimaryLen)
	}
	if res.ShadowLen != 1 {
		t.Errorf("expected ShadowLen=1, got %d", res.ShadowLen)
	}
}
