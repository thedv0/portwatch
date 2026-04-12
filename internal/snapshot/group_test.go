package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func gport(proto string, port, pid int, process string) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, PID: pid, Process: process}
}

func TestGroupPorts_ByProtocol(t *testing.T) {
	ports := []scanner.Port{
		gport("tcp", 80, 1, "nginx"),
		gport("udp", 53, 2, "dns"),
		gport("tcp", 443, 3, "nginx"),
	}
	groups := GroupPorts(ports, GroupByProtocol)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "tcp" || len(groups[0].Ports) != 2 {
		t.Errorf("unexpected tcp group: %+v", groups[0])
	}
	if groups[1].Key != "udp" || len(groups[1].Ports) != 1 {
		t.Errorf("unexpected udp group: %+v", groups[1])
	}
}

func TestGroupPorts_ByProcess(t *testing.T) {
	ports := []scanner.Port{
		gport("tcp", 80, 1, "nginx"),
		gport("tcp", 8080, 2, "caddy"),
		gport("tcp", 443, 3, "nginx"),
	}
	groups := GroupPorts(ports, GroupByProcess)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "caddy" || len(groups[0].Ports) != 1 {
		t.Errorf("unexpected caddy group: %+v", groups[0])
	}
	if groups[1].Key != "nginx" || len(groups[1].Ports) != 2 {
		t.Errorf("unexpected nginx group: %+v", groups[1])
	}
}

func TestGroupPorts_ByPID(t *testing.T) {
	ports := []scanner.Port{
		gport("tcp", 80, 10, "nginx"),
		gport("tcp", 443, 10, "nginx"),
		gport("udp", 53, 20, "dns"),
	}
	groups := GroupPorts(ports, GroupByPID)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupPorts_EmptyProtocolFallback(t *testing.T) {
	ports := []scanner.Port{
		gport("", 9999, 0, ""),
	}
	groups := GroupPorts(ports, GroupByProtocol)
	if len(groups) != 1 || groups[0].Key != "unknown" {
		t.Errorf("expected 'unknown' group, got %+v", groups)
	}
}

func TestGroupPorts_EmptyInput(t *testing.T) {
	groups := GroupPorts(nil, GroupByProtocol)
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestGroupPorts_SortedKeys(t *testing.T) {
	ports := []scanner.Port{
		gport("udp", 53, 1, "dns"),
		gport("tcp", 80, 2, "nginx"),
	}
	groups := GroupPorts(ports, GroupByProtocol)
	if groups[0].Key != "tcp" {
		t.Errorf("expected first key to be tcp, got %s", groups[0].Key)
	}
}
