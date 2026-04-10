package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, port uint16, pid int, process string) scanner.PortState {
	return scanner.PortState{
		Protocol: proto,
		     port,
		PID:      pid,
		Process:  process,
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	ports := []scanner.PortState{
		makePort("tcp", 80, 1, "nginx"),
		makePort("udp", 53, 2, "systemd"),
	}
	got := Filter(ports, FilterOptions{})
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestFilter_ByProtocol(t *testing.T) {
	ports := []scanner.PortState{
		makePort("tcp", 80, 1, "nginx"),
		makePort("udp", 53, 2, "systemd"),
		makePort("TCP", 443, 3, "nginx"),
	}
	got := Filter(ports, FilterOptions{Protocol: "tcp"})
	if len(got) != 2 {
		t.Fatalf("expected 2 tcp ports, got %d", len(got))
	}
}

func TestFilter_ByPortRange(t *testing.T) {
	ports := []scanner.PortState{
		makePort("tcp", 22, 1, "sshd"),
		makePort("tcp", 80, 2, "nginx"),
		makePort("tcp", 8080, 3, "app"),
	}
	got := Filter(ports, FilterOptions{MinPort: 79, MaxPort: 8079})
	if len(got) != 1 || got[0].Port != 80 {
		t.Fatalf("expected only port 80, got %v", got)
	}
}

func TestFilter_PIDZeroOnly(t *testing.T) {
	ports := []scanner.PortState{
		makePort("tcp", 80, 0, ""),
		makePort("tcp", 443, 5, "nginx"),
	}
	got := Filter(ports, FilterOptions{PIDZeroOnly: true})
	if len(got) != 1 || got[0].Port != 80 {
		t.Fatalf("expected only pid-zero port, got %v", got)
	}
}

func TestFilter_ByProcessName_CaseInsensitive(t *testing.T) {
	ports := []scanner.PortState{
		makePort("tcp", 80, 1, "Nginx"),
		makePort("tcp", 22, 2, "sshd"),
		makePort("tcp", 8080, 3, "nginx-proxy"),
	}
	got := Filter(ports, FilterOptions{ProcessName: "nginx"})
	if len(got) != 2 {
		t.Fatalf("expected 2 nginx ports, got %d", len(got))
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	got := Filter(nil, FilterOptions{Protocol: "tcp"})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestFilter_CombinedOptions(t *testing.T) {
	ports := []scanner.PortState{
		makePort("tcp", 80, 1, "nginx"),
		makePort("tcp", 443, 1, "nginx"),
		makePort("udp", 80, 2, "nginx"),
	}
	got := Filter(ports, FilterOptions{Protocol: "tcp", MinPort: 80, MaxPort: 80})
	if len(got) != 1 || got[0].Port != 80 || got[0].Protocol != "tcp" {
		t.Fatalf("expected exactly tcp:80, got %v", got)
	}
}
