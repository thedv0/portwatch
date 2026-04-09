package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func port(proto string, num int, pid int) scanner.Port {
	return scanner.Port{Protocol: proto, Port: num, PID: pid}
}

func TestDiff_NoChange(t *testing.T) {
	ports := []scanner.Port{port("tcp", 80, 1), port("tcp", 443, 2)}
	res := Diff(ports, ports)
	if len(res.Added) != 0 || len(res.Removed) != 0 {
		t.Errorf("expected no changes, got added=%d removed=%d", len(res.Added), len(res.Removed))
	}
}

func TestDiff_Added(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80, 1)}
	curr := []scanner.Port{port("tcp", 80, 1), port("tcp", 8080, 99)}
	res := Diff(prev, curr)
	if len(res.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(res.Added))
	}
	if res.Added[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", res.Added[0].Port)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removed, got %d", len(res.Removed))
	}
}

func TestDiff_Removed(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80, 1), port("udp", 53, 5)}
	curr := []scanner.Port{port("tcp", 80, 1)}
	res := Diff(prev, curr)
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(res.Removed))
	}
	if res.Removed[0].Port != 53 {
		t.Errorf("expected port 53, got %d", res.Removed[0].Port)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	curr := []scanner.Port{port("tcp", 22, 10), port("tcp", 80, 11)}
	res := Diff(nil, curr)
	if len(res.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(res.Added))
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(res.Removed))
	}
}

func TestDiff_ProtocolDistinct(t *testing.T) {
	prev := []scanner.Port{port("tcp", 53, 5)}
	curr := []scanner.Port{port("udp", 53, 5)}
	res := Diff(prev, curr)
	if len(res.Added) != 1 || len(res.Removed) != 1 {
		t.Errorf("expected 1 added and 1 removed for protocol change, got added=%d removed=%d",
			len(res.Added), len(res.Removed))
	}
}
