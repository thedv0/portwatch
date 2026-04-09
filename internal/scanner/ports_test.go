package scanner

import (
	"net"
	"testing"
	"time"
)

func TestParsePortRange_Single(t *testing.T) {
	ports, err := ParsePortRange("80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0] != 80 {
		t.Errorf("expected [80], got %v", ports)
	}
}

func TestParsePortRange_Multiple(t *testing.T) {
	ports, err := ParsePortRange("22,80,443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{22, 80, 443}
	for i, p := range expected {
		if ports[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, ports[i])
		}
	}
}

func TestParsePortRange_Range(t *testing.T) {
	ports, err := ParsePortRange("8000-8003")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 4 {
		t.Errorf("expected 4 ports, got %d", len(ports))
	}
}

func TestParsePortRange_Mixed(t *testing.T) {
	ports, err := ParsePortRange("22,8000-8001,443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 4 {
		t.Errorf("expected 4 ports, got %d", len(ports))
	}
}

func TestParsePortRange_Invalid(t *testing.T) {
	_, err := ParsePortRange("abc")
	if err == nil {
		t.Error("expected error for invalid port, got nil")
	}
}

func TestParsePortRange_InvalidRange(t *testing.T) {
	_, err := ParsePortRange("9000-8000")
	if err == nil {
		t.Error("expected error for inverted range, got nil")
	}
}

func TestScanTCP_DetectsOpenPort(t *testing.T) {
	// Start a temporary listener on an ephemeral port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	s := NewScanner(500 * time.Millisecond)
	entries, err := s.ScanTCP([]int{port})
	if err != nil {
		t.Fatalf("ScanTCP error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 open port, got %d", len(entries))
	}
	if entries[0].Port != port {
		t.Errorf("expected port %d, got %d", port, entries[0].Port)
	}
}

func TestScanTCP_ClosedPort(t *testing.T) {
	s := NewScanner(200 * time.Millisecond)
	// Port 1 is almost certainly not open in a test environment.
	entries, err := s.ScanTCP([]int{1})
	if err != nil {
		t.Fatalf("ScanTCP error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(entries))
	}
}
