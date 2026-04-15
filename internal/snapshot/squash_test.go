package snapshot

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func sqPort(proto string, port, pid int, process string) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, PID: pid, Process: process}
}

func makeSquashSnap(t time.Time, ports ...scanner.Port) Snapshot {
	return Snapshot{Timestamp: t, Ports: ports}
}

func TestSquash_EmptyInput(t *testing.T) {
	res, err := Squash(nil, DefaultSquashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Ports) != 0 {
		t.Errorf("expected 0 ports, got %d", len(res.Ports))
	}
	if res.InputSnaps != 0 {
		t.Errorf("expected InputSnaps=0, got %d", res.InputSnaps)
	}
}

func TestSquash_InvalidStrategy(t *testing.T) {
	snaps := []Snapshot{makeSquashSnap(time.Now(), sqPort("tcp", 80, 1, "nginx"))}
	_, err := Squash(snaps, SquashOptions{Strategy: "random"})
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func TestSquash_NoDuplicates(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeSquashSnap(now, sqPort("tcp", 80, 1, "nginx")),
		makeSquashSnap(now.Add(time.Second), sqPort("tcp", 443, 2, "nginx")),
	}
	res, err := Squash(snaps, DefaultSquashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(res.Ports))
	}
	if res.InputSnaps != 2 {
		t.Errorf("expected InputSnaps=2, got %d", res.InputSnaps)
	}
}

func TestSquash_StrategyLast_KeepsLatest(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeSquashSnap(now, sqPort("tcp", 80, 10, "old")),
		makeSquashSnap(now.Add(time.Second), sqPort("tcp", 80, 20, "new")),
	}
	res, _ := Squash(snaps, SquashOptions{Strategy: "last"})
	if len(res.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(res.Ports))
	}
	if res.Ports[0].Process != "new" {
		t.Errorf("expected process=new, got %s", res.Ports[0].Process)
	}
}

func TestSquash_StrategyFirst_KeepsEarliest(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeSquashSnap(now, sqPort("tcp", 80, 10, "old")),
		makeSquashSnap(now.Add(time.Second), sqPort("tcp", 80, 20, "new")),
	}
	res, _ := Squash(snaps, SquashOptions{Strategy: "first"})
	if len(res.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(res.Ports))
	}
	if res.Ports[0].Process != "old" {
		t.Errorf("expected process=old, got %s", res.Ports[0].Process)
	}
}

func TestSquash_LabelPreserved(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{makeSquashSnap(now, sqPort("tcp", 22, 1, "sshd"))}
	res, _ := Squash(snaps, SquashOptions{Strategy: "last", Label: "my-label"})
	if res.Label != "my-label" {
		t.Errorf("expected label=my-label, got %s", res.Label)
	}
}
