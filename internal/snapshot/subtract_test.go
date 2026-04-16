package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func subPort(port int, proto, process string, pid int) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func TestSubtract_EmptyRight_ReturnsAll(t *testing.T) {
	left := []scanner.Port{subPort(80, "tcp", "nginx", 1), subPort(443, "tcp", "nginx", 1)}
	got := Subtract(left, nil, DefaultSubtractOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestSubtract_DefaultKey_RemovesByProtoPort(t *testing.T) {
	left := []scanner.Port{
		subPort(80, "tcp", "nginx", 1),
		subPort(443, "tcp", "nginx", 1),
		subPort(8080, "tcp", "app", 2),
	}
	right := []scanner.Port{subPort(80, "tcp", "other", 99)}
	got := Subtract(left, right, DefaultSubtractOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
	for _, p := range got {
		if p.Port == 80 {
			t.Error("port 80 should have been subtracted")
		}
	}
}

func TestSubtract_ByPort_IgnoresProtocol(t *testing.T) {
	left := []scanner.Port{
		subPort(80, "tcp", "nginx", 1),
		subPort(80, "udp", "other", 2),
		subPort(443, "tcp", "nginx", 1),
	}
	right := []scanner.Port{subPort(80, "tcp", "nginx", 1)}
	opts := SubtractOptions{ByPort: true}
	got := Subtract(left, right, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got))
	}
	if got[0].Port != 443 {
		t.Errorf("expected port 443, got %d", got[0].Port)
	}
}

func TestSubtract_ByProcess_RemovesMatchingProcess(t *testing.T) {
	left := []scanner.Port{
		subPort(80, "tcp", "nginx", 1),
		subPort(8080, "tcp", "nginx", 2),
		subPort(443, "tcp", "app", 3),
	}
	right := []scanner.Port{subPort(9999, "udp", "nginx", 99)}
	opts := SubtractOptions{ByProcess: true}
	got := Subtract(left, right, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got))
	}
	if got[0].Process != "app" {
		t.Errorf("expected process 'app', got %q", got[0].Process)
	}
}

func TestSubtract_AllRemoved_ReturnsEmpty(t *testing.T) {
	left := []scanner.Port{subPort(80, "tcp", "nginx", 1)}
	right := []scanner.Port{subPort(80, "tcp", "nginx", 1)}
	got := Subtract(left, right, DefaultSubtractOptions())
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d ports", len(got))
	}
}

func TestSubtract_EmptyLeft_ReturnsEmpty(t *testing.T) {
	right := []scanner.Port{subPort(80, "tcp", "nginx", 1)}
	got := Subtract(nil, right, DefaultSubtractOptions())
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}
