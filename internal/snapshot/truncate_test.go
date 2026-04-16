package snapshot

import (
	"testing"
	"time"
)

func trPort(proto string, port int, seenAt time.Time) Port {
	return Port{Protocol: proto, Port: port, SeenAt: seenAt}
}

func makeTruncSnap(ports []Port) Snapshot {
	return Snapshot{Ports: ports, Timestamp: time.Now()}
}

func TestTruncate_NoOptions_ReturnsAll(t *testing.T) {
	snaps := []Snapshot{makeTruncSnap([]Port{
		trPort("tcp", 80, time.Time{}),
		trPort("tcp", 443, time.Time{}),
	})}
	out := Truncate(snaps, DefaultTruncateOptions())
	if len(out[0].Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out[0].Ports))
	}
}

func TestTruncate_MaxPorts_LimitsCount(t *testing.T) {
	snaps := []Snapshot{makeTruncSnap([]Port{
		trPort("tcp", 80, time.Time{}),
		trPort("tcp", 443, time.Time{}),
		trPort("tcp", 8080, time.Time{}),
	})}
	opts := DefaultTruncateOptions()
	opts.MaxPorts = 2
	out := Truncate(snaps, opts)
	if len(out[0].Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out[0].Ports))
	}
}

func TestTruncate_Before_RemovesOldPorts(t *testing.T) {
	now := time.Now()
	old := now.Add(-2 * time.Hour)
	snaps := []Snapshot{makeTruncSnap([]Port{
		trPort("tcp", 80, old),
		trPort("tcp", 443, now),
	})}
	opts := DefaultTruncateOptions()
	opts.Before = now.Add(-time.Hour)
	out := Truncate(snaps, opts)
	if len(out[0].Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out[0].Ports))
	}
	if out[0].Ports[0].Port != 443 {
		t.Errorf("expected port 443, got %d", out[0].Ports[0].Port)
	}
}

func TestTruncate_Protocol_OnlyFiltersMatchingProto(t *testing.T) {
	now := time.Now()
	old := now.Add(-2 * time.Hour)
	snaps := []Snapshot{makeTruncSnap([]Port{
		trPort("tcp", 80, old),
		trPort("udp", 53, old),
	})}
	opts := DefaultTruncateOptions()
	opts.Before = now.Add(-time.Hour)
	opts.Protocols = []string{"tcp"}
	out := Truncate(snaps, opts)
	// udp port should survive because protocol filter skips it
	if len(out[0].Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out[0].Ports))
	}
	if out[0].Ports[0].Protocol != "udp" {
		t.Errorf("expected udp port to survive, got %s", out[0].Ports[0].Protocol)
	}
}

func TestTruncate_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := Truncate(nil, DefaultTruncateOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %d", len(out))
	}
}

func TestTruncate_MaxPortsBeyondLength_ReturnsAll(t *testing.T) {
	snaps := []Snapshot{makeTruncSnap([]Port{
		trPort("tcp", 80, time.Time{}),
	})}
	opts := DefaultTruncateOptions()
	opts.MaxPorts = 100
	out := Truncate(snaps, opts)
	if len(out[0].Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out[0].Ports))
	}
}
