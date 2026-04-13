package snapshot

import (
	"testing"
	"time"
)

func prunePort(port int, proto string, seenAt time.Time) PortState {
	return PortState{Port: port, Protocol: proto, SeenAt: seenAt}
}

func TestPrune_NoOptions_ReturnsAll(t *testing.T) {
	ports := []PortState{
		prunePort(80, "tcp", time.Now()),
		prunePort(443, "tcp", time.Now()),
	}
	result := Prune(ports, DefaultPruneOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestPrune_Blacklist(t *testing.T) {
	ports := []PortState{
		prunePort(80, "tcp", time.Now()),
		prunePort(8080, "tcp", time.Now()),
		prunePort(443, "tcp", time.Now()),
	}
	opts := DefaultPruneOptions()
	opts.PortBlacklist = []int{8080}
	result := Prune(ports, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	for _, p := range result {
		if p.Port == 8080 {
			t.Error("blacklisted port 8080 should have been removed")
		}
	}
}

func TestPrune_AllowedProtocols(t *testing.T) {
	ports := []PortState{
		prunePort(53, "udp", time.Now()),
		prunePort(80, "tcp", time.Now()),
		prunePort(443, "tcp", time.Now()),
	}
	opts := DefaultPruneOptions()
	opts.AllowedProtocols = []string{"tcp"}
	result := Prune(ports, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 tcp ports, got %d", len(result))
	}
	for _, p := range result {
		if p.Protocol != "tcp" {
			t.Errorf("unexpected protocol %s", p.Protocol)
		}
	}
}

func TestPrune_MaxAge(t *testing.T) {
	old := time.Now().Add(-2 * time.Hour)
	recent := time.Now().Add(-1 * time.Minute)
	ports := []PortState{
		prunePort(80, "tcp", old),
		prunePort(443, "tcp", recent),
	}
	opts := DefaultPruneOptions()
	opts.MaxAge = 30 * time.Minute
	result := Prune(ports, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Port != 443 {
		t.Errorf("expected port 443, got %d", result[0].Port)
	}
}

func TestPrune_MaxPorts(t *testing.T) {
	ports := []PortState{
		prunePort(80, "tcp", time.Now()),
		prunePort(443, "tcp", time.Now()),
		prunePort(8080, "tcp", time.Now()),
	}
	opts := DefaultPruneOptions()
	opts.MaxPorts = 2
	result := Prune(ports, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestPrune_SortedByPort(t *testing.T) {
	ports := []PortState{
		prunePort(8080, "tcp", time.Now()),
		prunePort(80, "tcp", time.Now()),
		prunePort(443, "tcp", time.Now()),
	}
	result := Prune(ports, DefaultPruneOptions())
	if result[0].Port != 80 || result[1].Port != 443 || result[2].Port != 8080 {
		t.Error("result is not sorted by port number")
	}
}

func TestPrune_EmptyInput(t *testing.T) {
	result := Prune(nil, DefaultPruneOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}
