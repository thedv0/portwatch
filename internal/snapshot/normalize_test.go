package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func normPort(port int, proto, process string, pid int) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func TestNormalize_LowercaseProtocol(t *testing.T) {
	input := []scanner.Port{normPort(80, "TCP", "nginx", 100)}
	out := Normalize(input, NormalizeOptions{LowercaseProtocol: true})
	if out[0].Protocol != "tcp" {
		t.Errorf("expected tcp, got %s", out[0].Protocol)
	}
}

func TestNormalize_TrimProcessName(t *testing.T) {
	input := []scanner.Port{normPort(443, "tcp", "  sshd  ", 200)}
	out := Normalize(input, NormalizeOptions{TrimProcessName: true})
	if out[0].Process != "sshd" {
		t.Errorf("expected 'sshd', got %q", out[0].Process)
	}
}

func TestNormalize_ZeroInvalidPID(t *testing.T) {
	input := []scanner.Port{normPort(22, "tcp", "sshd", -5)}
	out := Normalize(input, NormalizeOptions{ZeroInvalidPID: true})
	if out[0].PID != 0 {
		t.Errorf("expected PID 0, got %d", out[0].PID)
	}
}

func TestNormalize_ClampPortHigh(t *testing.T) {
	input := []scanner.Port{normPort(99999, "tcp", "proc", 1)}
	out := Normalize(input, NormalizeOptions{ClampPort: true})
	if out[0].Port != 65535 {
		t.Errorf("expected 65535, got %d", out[0].Port)
	}
}

func TestNormalize_ClampPortLow(t *testing.T) {
	input := []scanner.Port{normPort(-10, "tcp", "proc", 1)}
	out := Normalize(input, NormalizeOptions{ClampPort: true})
	if out[0].Port != 0 {
		t.Errorf("expected 0, got %d", out[0].Port)
	}
}

func TestNormalize_DefaultOptions_AllApplied(t *testing.T) {
	input := []scanner.Port{normPort(-1, "UDP", "  app  ", -3)}
	out := Normalize(input, DefaultNormalizeOptions())
	p := out[0]
	if p.Protocol != "udp" {
		t.Errorf("protocol: expected udp, got %s", p.Protocol)
	}
	if p.Process != "app" {
		t.Errorf("process: expected 'app', got %q", p.Process)
	}
	if p.PID != 0 {
		t.Errorf("pid: expected 0, got %d", p.PID)
	}
	if p.Port != 0 {
		t.Errorf("port: expected 0, got %d", p.Port)
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	input := []scanner.Port{normPort(80, "TCP", "nginx", 10)}
	_ = Normalize(input, DefaultNormalizeOptions())
	if input[0].Protocol != "TCP" {
		t.Error("input slice was mutated")
	}
}

func TestNormalize_EmptyInput(t *testing.T) {
	out := Normalize(nil, DefaultNormalizeOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
