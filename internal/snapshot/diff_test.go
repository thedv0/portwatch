package snapshot_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func port(p int, proto, proc string) scanner.Port {
	return scanner.Port{Port: p, Protocol: proto, Process: proc}
}

func TestDiff_NoChange(t *testing.T) {
	ports := []scanner.Port{port(80, "tcp", "nginx"), port(443, "tcp", "nginx")}
	changes := snapshot.Diff(ports, ports)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestDiff_Added(t *testing.T) {
	prev := []scanner.Port{port(80, "tcp", "nginx")}
	curr := []scanner.Port{port(80, "tcp", "nginx"), port(8080, "tcp", "python")}

	changes := snapshot.Diff(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != snapshot.Added {
		t.Errorf("expected Added, got %s", changes[0].Kind)
	}
	if changes[0].Port.Port != 8080 {
		t.Errorf("expected port 8080, got %d", changes[0].Port.Port)
	}
}

func TestDiff_Removed(t *testing.T) {
	prev := []scanner.Port{port(80, "tcp", "nginx"), port(22, "tcp", "sshd")}
	curr := []scanner.Port{port(80, "tcp", "nginx")}

	changes := snapshot.Diff(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != snapshot.Removed {
		t.Errorf("expected Removed, got %s", changes[0].Kind)
	}
	if changes[0].Port.Port != 22 {
		t.Errorf("expected port 22, got %d", changes[0].Port.Port)
	}
}

func TestDiff_EmptyPrev(t *testing.T) {
	curr := []scanner.Port{port(80, "tcp", "nginx"), port(443, "tcp", "nginx")}
	changes := snapshot.Diff(nil, curr)
	if len(changes) != 2 {
		t.Errorf("expected 2 added changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.Kind != snapshot.Added {
			t.Errorf("expected Added, got %s", c.Kind)
		}
	}
}

func TestDiff_ProtocolDistinct(t *testing.T) {
	prev := []scanner.Port{port(53, "tcp", "named")}
	curr := []scanner.Port{port(53, "udp", "named")}

	changes := snapshot.Diff(prev, curr)
	if len(changes) != 2 {
		t.Errorf("expected 2 changes (tcp removed, udp added), got %d", len(changes))
	}
}
