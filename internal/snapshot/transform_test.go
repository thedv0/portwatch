package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func txPort(port int, proto, process string, pid int) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func TestTransform_NoOp(t *testing.T) {
	ports := []scanner.Port{txPort(80, "tcp", "nginx", 100)}
	opts := DefaultTransformOptions()
	out := Transform(ports, opts)
	if out[0].Port != 80 || out[0].Protocol != "tcp" || out[0].Process != "nginx" {
		t.Fatalf("unexpected change: %+v", out[0])
	}
}

func TestTransform_ForceUpperProtocol(t *testing.T) {
	ports := []scanner.Port{txPort(443, "tcp", "nginx", 1)}
	opts := DefaultTransformOptions()
	opts.ForceUpperProtocol = true
	out := Transform(ports, opts)
	if out[0].Protocol != "TCP" {
		t.Fatalf("expected TCP, got %s", out[0].Protocol)
	}
}

func TestTransform_MapProtocol(t *testing.T) {
	ports := []scanner.Port{txPort(53, "udp", "dns", 2)}
	opts := DefaultTransformOptions()
	opts.MapProtocol = map[string]string{"udp": "UDP4"}
	out := Transform(ports, opts)
	if out[0].Protocol != "UDP4" {
		t.Fatalf("expected UDP4, got %s", out[0].Protocol)
	}
}

func TestTransform_MapProtocol_TakesPrecedenceOverUpper(t *testing.T) {
	ports := []scanner.Port{txPort(53, "udp", "dns", 2)}
	opts := DefaultTransformOptions()
	opts.MapProtocol = map[string]string{"udp": "datagram"}
	opts.ForceUpperProtocol = true
	out := Transform(ports, opts)
	// MapProtocol runs first; ForceUpperProtocol then uppercases the result.
	if out[0].Protocol != "DATAGRAM" {
		t.Fatalf("expected DATAGRAM, got %s", out[0].Protocol)
	}
}

func TestTransform_RenameProcess(t *testing.T) {
	ports := []scanner.Port{txPort(8080, "tcp", "java", 50)}
	opts := DefaultTransformOptions()
	opts.RenameProcess = map[string]string{"java": "app-server"}
	out := Transform(ports, opts)
	if out[0].Process != "app-server" {
		t.Fatalf("expected app-server, got %s", out[0].Process)
	}
}

func TestTransform_OffsetPorts(t *testing.T) {
	ports := []scanner.Port{txPort(1000, "tcp", "svc", 10)}
	opts := DefaultTransformOptions()
	opts.OffsetPorts = 100
	out := Transform(ports, opts)
	if out[0].Port != 1100 {
		t.Fatalf("expected 1100, got %d", out[0].Port)
	}
}

func TestTransform_OffsetPorts_ClampsLow(t *testing.T) {
	ports := []scanner.Port{txPort(1, "tcp", "svc", 10)}
	opts := DefaultTransformOptions()
	opts.OffsetPorts = -500
	out := Transform(ports, opts)
	if out[0].Port != 1 {
		t.Fatalf("expected clamped to 1, got %d", out[0].Port)
	}
}

func TestTransform_OffsetPorts_ClampsHigh(t *testing.T) {
	ports := []scanner.Port{txPort(65500, "tcp", "svc", 10)}
	opts := DefaultTransformOptions()
	opts.OffsetPorts = 1000
	out := Transform(ports, opts)
	if out[0].Port != 65535 {
		t.Fatalf("expected clamped to 65535, got %d", out[0].Port)
	}
}

func TestTransform_DoesNotMutateOriginal(t *testing.T) {
	original := []scanner.Port{txPort(80, "tcp", "nginx", 1)}
	opts := DefaultTransformOptions()
	opts.ForceUpperProtocol = true
	opts.OffsetPorts = 10
	_ = Transform(original, opts)
	if original[0].Protocol != "tcp" || original[0].Port != 80 {
		t.Fatal("original slice was mutated")
	}
}
