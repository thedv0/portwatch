package snapshot

import (
	"testing"
	"time"
)

var fixedEvictNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func evPort(port int, proto string) PortState {
	return PortState{Port: port, Protocol: proto}
}

func TestEvict_ByAge_RemovesOldPorts(t *testing.T) {
	ports := []PortState{evPort(80, "tcp"), evPort(443, "tcp"), evPort(8080, "tcp")}
	lastSeen := map[string]time.Time{
		"tcp:80":   fixedEvictNow.Add(-60 * time.Minute), // old
		"tcp:443":  fixedEvictNow.Add(-5 * time.Minute),  // recent
		"tcp:8080": fixedEvictNow.Add(-10 * time.Minute), // recent
	}
	opts := DefaultEvictOptions()
	opts.MaxAge = 30 * time.Minute
	opts.Now = func() time.Time { return fixedEvictNow }

	result, err := Evict(ports, lastSeen, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
}

func TestEvict_ByAge_KeepsUnseenPorts(t *testing.T) {
	ports := []PortState{evPort(22, "tcp")}
	lastSeen := map[string]time.Time{} // not in map
	opts := DefaultEvictOptions()
	opts.Now = func() time.Time { return fixedEvictNow }

	result, err := Evict(ports, lastSeen, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 port, got %d", len(result))
	}
}

func TestEvict_ByCount_KeepsNewest(t *testing.T) {
	ports := []PortState{evPort(80, "tcp"), evPort(443, "tcp"), evPort(8080, "tcp")}
	lastSeen := map[string]time.Time{
		"tcp:80":   fixedEvictNow.Add(-30 * time.Minute),
		"tcp:443":  fixedEvictNow.Add(-2 * time.Minute),
		"tcp:8080": fixedEvictNow.Add(-1 * time.Minute),
	}
	opts := DefaultEvictOptions()
	opts.Policy = EvictByCount
	opts.MaxCount = 2
	opts.Now = func() time.Time { return fixedEvictNow }

	result, err := Evict(ports, lastSeen, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
}

func TestEvict_ByCount_InvalidMaxCount(t *testing.T) {
	opts := DefaultEvictOptions()
	opts.Policy = EvictByCount
	opts.MaxCount = 0

	_, err := Evict([]PortState{evPort(80, "tcp")}, map[string]time.Time{}, opts)
	if err == nil {
		t.Fatal("expected error for MaxCount=0")
	}
}

func TestEvict_ByIdleTime_RemovesIdlePorts(t *testing.T) {
	ports := []PortState{evPort(9000, "udp"), evPort(9001, "udp")}
	lastSeen := map[string]time.Time{
		"udp:9000": fixedEvictNow.Add(-20 * time.Minute), // idle
		"udp:9001": fixedEvictNow.Add(-3 * time.Minute),  // active
	}
	opts := DefaultEvictOptions()
	opts.Policy = EvictByIdleTime
	opts.IdleTime = 10 * time.Minute
	opts.Now = func() time.Time { return fixedEvictNow }

	result, err := Evict(ports, lastSeen, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Port != 9001 {
		t.Fatalf("expected only port 9001, got %+v", result)
	}
}

func TestEvict_UnknownPolicy_ReturnsError(t *testing.T) {
	opts := DefaultEvictOptions()
	opts.Policy = EvictPolicy(99)

	_, err := Evict([]PortState{}, map[string]time.Time{}, opts)
	if err == nil {
		t.Fatal("expected error for unknown policy")
	}
}

func TestDefaultEvictOptions_Values(t *testing.T) {
	o := DefaultEvictOptions()
	if o.MaxAge <= 0 {
		t.Error("MaxAge should be positive")
	}
	if o.MaxCount <= 0 {
		t.Error("MaxCount should be positive")
	}
	if o.IdleTime <= 0 {
		t.Error("IdleTime should be positive")
	}
}
