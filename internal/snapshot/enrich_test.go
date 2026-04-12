package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func enrichPort(port int, proto, process string, pid int) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func TestEnrich_NormalizesProtocol(t *testing.T) {
	ports := []scanner.Port{enrichPort(80, "TCP", "nginx", 100)}
	out := Enrich(ports, DefaultEnrichOptions())
	if out[0].Protocol != "tcp" {
		t.Errorf("expected protocol 'tcp', got %q", out[0].Protocol)
	}
}

func TestEnrich_ResolvesWellKnown(t *testing.T) {
	ports := []scanner.Port{enrichPort(22, "tcp", "", 0)}
	out := Enrich(ports, DefaultEnrichOptions())
	if out[0].Process != "ssh" {
		t.Errorf("expected process 'ssh', got %q", out[0].Process)
	}
}

func TestEnrich_DoesNotOverwriteExistingProcess(t *testing.T) {
	ports := []scanner.Port{enrichPort(22, "tcp", "custom-sshd", 42)}
	out := Enrich(ports, DefaultEnrichOptions())
	if out[0].Process != "custom-sshd" {
		t.Errorf("expected process 'custom-sshd', got %q", out[0].Process)
	}
}

func TestEnrich_UnknownPortNoChange(t *testing.T) {
	ports := []scanner.Port{enrichPort(9999, "udp", "", 0)}
	out := Enrich(ports, DefaultEnrichOptions())
	if out[0].Process != "" {
		t.Errorf("expected empty process, got %q", out[0].Process)
	}
}

func TestEnrich_SkipNormalize(t *testing.T) {
	opts := DefaultEnrichOptions()
	opts.NormalizeProtocol = false
	ports := []scanner.Port{enrichPort(80, "TCP", "", 0)}
	out := Enrich(ports, opts)
	if out[0].Protocol != "TCP" {
		t.Errorf("expected protocol 'TCP' unchanged, got %q", out[0].Protocol)
	}
}

func TestEnrich_DoesNotMutateInput(t *testing.T) {
	original := []scanner.Port{enrichPort(22, "TCP", "", 0)}
	_ = Enrich(original, DefaultEnrichOptions())
	if original[0].Protocol != "TCP" {
		t.Error("input slice was mutated")
	}
}

func TestEnrich_EmptyInput(t *testing.T) {
	out := Enrich(nil, DefaultEnrichOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d items", len(out))
	}
}
