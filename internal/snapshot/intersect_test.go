package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func ixPort(port int, proto, process string, pid int) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func TestIntersect_EmptyInput(t *testing.T) {
	result := Intersect(nil, DefaultIntersectOptions())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestIntersect_SingleSnapshot_ReturnsAllPorts(t *testing.T) {
	snaps := []Snapshot{
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 100), ixPort(443, "tcp", "nginx", 100)}},
	}
	result := Intersect(snaps, DefaultIntersectOptions())
	if len(result) != 2 {
		t.Errorf("expected 2 ports, got %d", len(result))
	}
}

func TestIntersect_CommonPortsOnly(t *testing.T) {
	snaps := []Snapshot{
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1), ixPort(443, "tcp", "nginx", 1), ixPort(8080, "tcp", "app", 2)}},
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1), ixPort(9090, "tcp", "prom", 3)}},
	}
	result := Intersect(snaps, DefaultIntersectOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 port, got %d", len(result))
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80, got %d", result[0].Port)
	}
}

func TestIntersect_NoCommonPorts_ReturnsEmpty(t *testing.T) {
	snaps := []Snapshot{
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1)}},
		{Ports: []scanner.Port{ixPort(443, "tcp", "nginx", 1)}},
	}
	result := Intersect(snaps, DefaultIntersectOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestIntersect_ThreeSnapshots_AllMustMatch(t *testing.T) {
	snaps := []Snapshot{
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1), ixPort(443, "tcp", "nginx", 1)}},
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1), ixPort(443, "tcp", "nginx", 1)}},
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1)}},
	}
	result := Intersect(snaps, DefaultIntersectOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 port, got %d", len(result))
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80, got %d", result[0].Port)
	}
}

func TestIntersect_ByProcess_KeyField(t *testing.T) {
	opts := IntersectOptions{KeyFields: []string{"process"}}
	snaps := []Snapshot{
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1), ixPort(9000, "udp", "other", 2)}},
		{Ports: []scanner.Port{ixPort(8080, "tcp", "nginx", 5)}},
	}
	result := Intersect(snaps, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Process != "nginx" {
		t.Errorf("expected nginx, got %s", result[0].Process)
	}
}

func TestIntersect_DefaultOptions_UsedWhenFieldsEmpty(t *testing.T) {
	opts := IntersectOptions{} // empty KeyFields triggers default
	snaps := []Snapshot{
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1)}},
		{Ports: []scanner.Port{ixPort(80, "tcp", "nginx", 1)}},
	}
	result := Intersect(snaps, opts)
	if len(result) != 1 {
		t.Errorf("expected 1 port, got %d", len(result))
	}
}
