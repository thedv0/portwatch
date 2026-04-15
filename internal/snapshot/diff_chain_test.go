package snapshot

import (
	"testing"
	"time"
)

func makeChainSnap(t time.Time, ports []Port) Snapshot {
	return Snapshot{Timestamp: t, Ports: ports}
}

func chainPort(proto string, num int, pid int) Port {
	return Port{Protocol: proto, Port: num, PID: pid, Process: "proc"}
}

func TestBuildDiffChain_EmptyInput(t *testing.T) {
	_, err := BuildDiffChain(nil, DefaultDiffChainOptions())
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestBuildDiffChain_SingleSnapshot(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{makeChainSnap(now, []Port{chainPort("tcp", 80, 1)})}
	chain, err := BuildDiffChain(snaps, DefaultDiffChainOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(chain))
	}
	if len(chain[0].Diff.Added) != 0 || len(chain[0].Diff.Removed) != 0 {
		t.Error("first entry should have empty diff")
	}
}

func TestBuildDiffChain_DiffPopulated(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeChainSnap(now, []Port{chainPort("tcp", 80, 1)}),
		makeChainSnap(now.Add(time.Minute), []Port{chainPort("tcp", 80, 1), chainPort("tcp", 443, 2)}),
	}
	chain, err := BuildDiffChain(snaps, DefaultDiffChainOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain[1].Diff.Added) != 1 {
		t.Errorf("expected 1 added port, got %d", len(chain[1].Diff.Added))
	}
}

func TestBuildDiffChain_MaxEntries(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeChainSnap(now, nil),
		makeChainSnap(now.Add(time.Minute), nil),
		makeChainSnap(now.Add(2*time.Minute), nil),
	}
	opts := DefaultDiffChainOptions()
	opts.MaxEntries = 2
	chain, err := BuildDiffChain(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain) != 2 {
		t.Errorf("expected 2 entries, got %d", len(chain))
	}
}

func TestBuildDiffChain_SinceFilter(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeChainSnap(now.Add(-2*time.Hour), nil),
		makeChainSnap(now.Add(-30*time.Minute), nil),
		makeChainSnap(now, nil),
	}
	opts := DefaultDiffChainOptions()
	opts.Since = now.Add(-time.Hour)
	chain, err := BuildDiffChain(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chain) != 2 {
		t.Errorf("expected 2 entries after Since filter, got %d", len(chain))
	}
}

func TestBuildDiffChain_IndexField(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeChainSnap(now, nil),
		makeChainSnap(now.Add(time.Minute), nil),
	}
	chain, _ := BuildDiffChain(snaps, DefaultDiffChainOptions())
	for i, e := range chain {
		if e.Index != i {
			t.Errorf("entry %d has Index %d", i, e.Index)
		}
	}
}
