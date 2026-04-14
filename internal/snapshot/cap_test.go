package snapshot

import (
	"testing"
)

func capPort(port int, pid int) PortState {
	return PortState{Port: port, PID: pid, Protocol: "tcp", Process: "test"}
}

func TestCap_NoLimit_ReturnsAll(t *testing.T) {
	ports := []PortState{capPort(80, 1), capPort(443, 2), capPort(22, 3)}
	opts := DefaultCapOptions()
	opts.MaxPorts = 0
	out := Cap(ports, opts)
	if len(out) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(out))
	}
}

func TestCap_LimitsToMaxPorts(t *testing.T) {
	ports := []PortState{capPort(80, 1), capPort(443, 2), capPort(8080, 3), capPort(9090, 4)}
	opts := DefaultCapOptions()
	opts.MaxPorts = 2
	out := Cap(ports, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out))
	}
}

func TestCap_SortByPort_RetainsLowest(t *testing.T) {
	ports := []PortState{capPort(9090, 1), capPort(22, 2), capPort(443, 3)}
	opts := DefaultCapOptions()
	opts.SortByPort = true
	opts.MaxPorts = 2
	out := Cap(ports, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out))
	}
	if out[0].Port != 22 {
		t.Errorf("expected first port 22, got %d", out[0].Port)
	}
	if out[1].Port != 443 {
		t.Errorf("expected second port 443, got %d", out[1].Port)
	}
}

func TestCap_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := Cap(nil, DefaultCapOptions())
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestCap_MaxGreaterThanLen_ReturnsAll(t *testing.T) {
	ports := []PortState{capPort(80, 1), capPort(443, 2)}
	opts := DefaultCapOptions()
	opts.MaxPorts = 100
	out := Cap(ports, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out))
	}
}

func TestCap_DoesNotMutateInput(t *testing.T) {
	ports := []PortState{capPort(9090, 1), capPort(22, 2), capPort(443, 3)}
	orig := make([]PortState, len(ports))
	copy(orig, ports)
	opts := DefaultCapOptions()
	opts.SortByPort = true
	opts.MaxPorts = 2
	Cap(ports, opts)
	for i, p := range ports {
		if p.Port != orig[i].Port {
			t.Errorf("input mutated at index %d: got %d want %d", i, p.Port, orig[i].Port)
		}
	}
}

func TestDefaultCapOptions_Values(t *testing.T) {
	opts := DefaultCapOptions()
	if opts.MaxPorts != 0 {
		t.Errorf("expected MaxPorts 0, got %d", opts.MaxPorts)
	}
	if !opts.SortByPort {
		t.Error("expected SortByPort true")
	}
}
