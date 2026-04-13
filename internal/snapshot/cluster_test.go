package snapshot

import (
	"testing"

	"github.com/netwatch/portwatch/internal/scanner"
)

func clport(port int, proto, process string, pid int) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func makeClusterSnaps(ports ...scanner.Port) []scanner.Snapshot {
	return []scanner.Snapshot{{Ports: ports}}
}

func TestCluster_ByProcess_GroupsSameProcess(t *testing.T) {
	snaps := makeClusterSnaps(
		clport(80, "tcp", "nginx", 100),
		clport(443, "tcp", "nginx", 100),
		clport(5432, "tcp", "postgres", 200),
	)
	result := Cluster(snaps, DefaultClusterOptions())
	if len(result["nginx"]) != 2 {
		t.Errorf("expected 2 nginx ports, got %d", len(result["nginx"]))
	}
	if len(result["postgres"]) != 1 {
		t.Errorf("expected 1 postgres port, got %d", len(result["postgres"]))
	}
}

func TestCluster_ByProtocol_GroupsCorrectly(t *testing.T) {
	snaps := makeClusterSnaps(
		clport(80, "tcp", "nginx", 1),
		clport(53, "udp", "dns", 2),
		clport(443, "tcp", "nginx", 1),
	)
	opts := ClusterOptions{By: "protocol", MinSize: 1}
	result := Cluster(snaps, opts)
	if len(result["tcp"]) != 2 {
		t.Errorf("expected 2 tcp ports, got %d", len(result["tcp"]))
	}
	if len(result["udp"]) != 1 {
		t.Errorf("expected 1 udp port, got %d", len(result["udp"]))
	}
}

func TestCluster_MinSize_FiltersSmallClusters(t *testing.T) {
	snaps := makeClusterSnaps(
		clport(80, "tcp", "nginx", 1),
		clport(443, "tcp", "nginx", 1),
		clport(5432, "tcp", "postgres", 2),
	)
	opts := ClusterOptions{By: "process", MinSize: 2}
	result := Cluster(snaps, opts)
	if _, ok := result["postgres"]; ok {
		t.Error("expected postgres cluster to be filtered out")
	}
	if len(result["nginx"]) != 2 {
		t.Errorf("expected 2 nginx ports, got %d", len(result["nginx"]))
	}
}

func TestCluster_EmptyInput_ReturnsEmpty(t *testing.T) {
	result := Cluster(nil, DefaultClusterOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d clusters", len(result))
	}
}

func TestCluster_UnknownProcess_FallsBack(t *testing.T) {
	snaps := makeClusterSnaps(clport(9999, "tcp", "", 0))
	result := Cluster(snaps, DefaultClusterOptions())
	if _, ok := result["unknown"]; !ok {
		t.Error("expected 'unknown' cluster for empty process name")
	}
}

func TestCluster_ByPort_KeyIsPortNumber(t *testing.T) {
	snaps := makeClusterSnaps(
		clport(80, "tcp", "nginx", 1),
		clport(80, "udp", "other", 2),
	)
	opts := ClusterOptions{By: "port", MinSize: 1}
	result := Cluster(snaps, opts)
	if len(result["80"]) != 2 {
		t.Errorf("expected 2 entries for port 80, got %d", len(result["80"]))
	}
}
