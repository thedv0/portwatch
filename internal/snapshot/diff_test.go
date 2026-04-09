package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func port(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Port: num}
}

func TestDiff_NoChange(t *testing.T) {
	ports := []scanner.Port{port("tcp", 80), port("tcp", 443)}
	res := Diff(ports, ports)
	if len(res.Added) != 0 || len(res.Removed) != 0 {
		t.Fatalf("expected no diff, got added=%v removed=%v", res.Added, res.Removed)
	}
}

func TestDiff_Added(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80)}
	curr := []scanner.Port{port("tcp", 80), port("tcp", 8080)}
	res := Diff(prev, curr)
	if len(res.Added) != 1 || res.Added[0].Port != 8080 {
		t.Fatalf("expected port 8080 added, got %v", res.Added)
	}
	if len(res.Removed) != 0 {
		t.Fatalf("expected no removals, got %v", res.Removed)
	}
}

func TestDiff_Removed(t *testing.T) {
	prev := []scanner.Port{port("tcp", 80), port("tcp", 9000)}
	curr := []scanner.Port{port("tcp", 80)}
	res := Diff(prev, curr)
	if len(res.Removed) != 1 || res.Removed[0].Port != 9000 {
		t.Fatalf("expected port 9000 removed, got %v", res.Removed)
	}
	if len(res.Added) != 0 {
		t.Fatalf("expected no additions, got %v", res.Added)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	curr := []scanner.Port{port("tcp", 22), port("udp", 53)}
	res := Diff(nil, curr)
	if len(res.Added) != 2 {
		t.Fatalf("expected 2 added, got %v", res.Added)
	}
}

func TestDiff_ProtocolDistinct(t *testing.T) {
	prev := []scanner.Port{port("tcp", 53)}
	curr := []scanner.Port{port("udp", 53)}
	res := Diff(prev, curr)
	if len(res.Added) != 1 || res.Added[0].Protocol != "udp" {
		t.Fatalf("expected udp:53 added, got %v", res.Added)
	}
	if len(res.Removed) != 1 || res.Removed[0].Protocol != "tcp" {
		t.Fatalf("expected tcp:53 removed, got %v", res.Removed)
	}
}
