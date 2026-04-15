package snapshot

import (
	"testing"
)

func spPort(port int, proto, process string) Port {
	return Port{Port: port, Protocol: proto, Process: process, PID: 100}
}

func TestSplit_EvenCount(t *testing.T) {
	ports := []Port{spPort(80, "tcp", "nginx"), spPort(443, "tcp", "nginx"), spPort(8080, "tcp", "app"), spPort(9090, "udp", "app")}
	opts := DefaultSplitOptions()
	opts.NumParts = 2
	opts.Field = "port"
	res, err := Split(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(res))
	}
	if len(res[0].Ports) != 2 || len(res[1].Ports) != 2 {
		t.Errorf("expected 2 ports per part")
	}
}

func TestSplit_ByProtocol(t *testing.T) {
	ports := []Port{spPort(80, "tcp", "nginx"), spPort(53, "udp", "dns"), spPort(443, "tcp", "nginx")}
	opts := DefaultSplitOptions()
	opts.Field = "protocol"
	res, err := Split(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 groups (tcp, udp), got %d", len(res))
	}
	// first group should be tcp (2 ports)
	if len(res[0].Ports) != 2 {
		t.Errorf("expected 2 tcp ports, got %d", len(res[0].Ports))
	}
}

func TestSplit_ByProcess(t *testing.T) {
	ports := []Port{spPort(80, "tcp", "nginx"), spPort(443, "tcp", "nginx"), spPort(3000, "tcp", "node")}
	opts := DefaultSplitOptions()
	opts.Field = "process"
	res, err := Split(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 process groups, got %d", len(res))
	}
}

func TestSplit_InvalidNumParts(t *testing.T) {
	ports := []Port{spPort(80, "tcp", "nginx")}
	opts := DefaultSplitOptions()
	opts.NumParts = 0
	_, err := Split(ports, opts)
	if err == nil {
		t.Error("expected error for NumParts=0")
	}
}

func TestSplit_UnknownField(t *testing.T) {
	ports := []Port{spPort(80, "tcp", "nginx")}
	opts := DefaultSplitOptions()
	opts.Field = "invalid"
	_, err := Split(ports, opts)
	if err == nil {
		t.Error("expected error for unknown field")
	}
}

func TestSplit_EmptyPorts(t *testing.T) {
	opts := DefaultSplitOptions()
	res, err := Split([]Port{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range res {
		if len(r.Ports) != 0 {
			t.Errorf("expected empty parts for empty input")
		}
	}
}
