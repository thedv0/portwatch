package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func port(proto string, number int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: number}
}

func TestDiff_NoChange(t *testing.T) {
	ports := []scanner.Port{port("tcp", 80), port("tcp", 443)}
	result := Diff(ports, ports)
	if len(result.Added) != 0 || len(result.Removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", result.Added, result.Removed)
	}
}

func TestDiff_Added(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80)}
	curr := []scanner.Port{port("tcp", 80), port("tcp", 8080)}
	result := Diff(prev, curr)
	if len(result.Added) != 1 || result.Added[0].Number != 8080 {
		t.Errorf("expected port 8080 added, got %v", result.Added)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removals, got %v", result.Removed)
	}
}

func TestDiff_Removed(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80), port("tcp", 9090)}
	curr := []scanner.Port{port("tcp", 80)}
	result := Diff(prev, curr)
	if len(result.Removed) != 1 || result.Removed[0].Number != 9090 {
		t.Errorf("expected port 9090 removed, got %v", result.Removed)
	}
	if len(result.Added) != 0 {
		t.Errorf("expected no additions, got %v", result.Added)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	curr := []scanner.Port{port("tcp", 22), port("udp", 53)}
	result := Diff(nil, curr)
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added ports, got %d", len(result.Added))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removals, got %v", result.Removed)
	}
}

func TestDiff_ProtocolDistinct(t *testing.T) {
	prev := []scanner.Port{port("tcp", 53)}
	curr := []scanner.Port{port("udp", 53)}
	result := Diff(prev, curr)
	if len(result.Added) != 1 || result.Added[0].Protocol != "udp" {
		t.Errorf("expected udp:53 added, got %v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed[0].Protocol != "tcp" {
		t.Errorf("expected tcp:53 removed, got %v", result.Removed)
	}
}
