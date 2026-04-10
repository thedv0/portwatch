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
	result := Diff(ports, ports)
	if len(result.Added) != 0 || len(result.Removed) != 0 {
		t.Errorf("expected no changes, got added=%d removed=%d", len(result.Added), len(result.Removed))
	}
}

func TestDiff_Added(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80, 1)}
	curr := []scanner.Port{port("tcp", 80, 1), port("tcp", 8080, 99)}
	result := Diff(prev, curr)
	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(result.Added))
	}
	if result.Added[0].Port != 8080 {
		t.Errorf("expected added port 8080, got %d", result.Added[0].Port)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removed, got %d", len(result.Removed))
	}
}

func TestDiff_Removed(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80, 1), port("udp", 53, 5)}
	curr := []scanner.Port{port("tcp", 80, 1)}
	result := Diff(prev, curr)
	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(result.Removed))
	}
	if result.Removed[0].Port != 53 {
		t.Errorf("expected removed port 53, got %d", result.Removed[0].Port)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	curr := []scanner.Port{port("tcp", 22, 10), port("tcp", 80, 11)}
	result := Diff(nil, curr)
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
}

func TestDiff_ProtocolDistinct(t *testing.T) {
	prev := []scanner.Port{port("tcp", 53, 5)}
	curr := []scanner.Port{port("udp", 53, 5)}
	result := Diff(prev, curr)
	if len(result.Added) != 1 || len(result.Removed) != 1 {
		t.Errorf("tcp and udp on same port should be distinct: added=%d removed=%d",
			len(result.Added), len(result.Removed))
	}
}
