package snapshot

import (
	"testing"
	"time"
)

func makeSnaps(sets [][]Port) []Snapshot {
	base := time.Now()
	snaps := make([]Snapshot, len(sets))
	for i, ports := range sets {
		snaps[i] = Snapshot{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Ports:     ports,
		}
	}
	return snaps
}

func TestAggregate_Empty(t *testing.T) {
	stats := Aggregate(nil)
	if stats.SampleCount != 0 {
		t.Errorf("expected 0 samples, got %d", stats.SampleCount)
	}
}

func TestAggregate_SampleCount(t *testing.T) {
	snaps := makeSnaps([][]Port{
		{{Port: 80, Protocol: "tcp"}},
		{{Port: 443, Protocol: "tcp"}},
	})
	stats := Aggregate(snaps)
	if stats.SampleCount != 2 {
		t.Errorf("expected 2 samples, got %d", stats.SampleCount)
	}
}

func TestAggregate_AvgOpen(t *testing.T) {
	snaps := makeSnaps([][]Port{
		{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}},
		{{Port: 22, Protocol: "tcp"}},
	})
	stats := Aggregate(snaps)
	// avg = (2+1)/2 = 1.5
	if stats.AvgOpen != 1.5 {
		t.Errorf("expected avg 1.5, got %f", stats.AvgOpen)
	}
}

func TestAggregate_MaxMin(t *testing.T) {
	snaps := makeSnaps([][]Port{
		{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}, {Port: 8080, Protocol: "tcp"}},
		{{Port: 22, Protocol: "tcp"}},
	})
	stats := Aggregate(snaps)
	if stats.MaxOpen != 3 {
		t.Errorf("expected max 3, got %d", stats.MaxOpen)
	}
	if stats.MinOpen != 1 {
		t.Errorf("expected min 1, got %d", stats.MinOpen)
	}
}

func TestAggregate_TopPorts(t *testing.T) {
	snaps := makeSnaps([][]Port{
		{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}},
		{{Port: 80, Protocol: "tcp"}, {Port: 22, Protocol: "tcp"}},
		{{Port: 80, Protocol: "tcp"}},
	})
	stats := Aggregate(snaps)
	if len(stats.TopPorts) == 0 {
		t.Fatal("expected top ports, got none")
	}
	if stats.TopPorts[0].Port != 80 {
		t.Errorf("expected port 80 at top, got %d", stats.TopPorts[0].Port)
	}
	if stats.TopPorts[0].Count != 3 {
		t.Errorf("expected count 3, got %d", stats.TopPorts[0].Count)
	}
}

func TestAggregate_TopProcesses(t *testing.T) {
	snaps := makeSnaps([][]Port{
		{{Port: 80, Protocol: "tcp", Process: "nginx"}, {Port: 443, Protocol: "tcp", Process: "nginx"}},
		{{Port: 22, Protocol: "tcp", Process: "sshd"}},
	})
	stats := Aggregate(snaps)
	if len(stats.TopProcesses) == 0 {
		t.Fatal("expected top processes, got none")
	}
	if stats.TopProcesses[0].Process != "nginx" {
		t.Errorf("expected nginx at top, got %s", stats.TopProcesses[0].Process)
	}
}

func TestAggregate_TimeRange(t *testing.T) {
	snaps := makeSnaps([][]Port{
		{{Port: 80, Protocol: "tcp"}},
		{{Port: 443, Protocol: "tcp"}},
	})
	stats := Aggregate(snaps)
	if !stats.To.After(stats.From) {
		t.Errorf("expected To after From, got From=%v To=%v", stats.From, stats.To)
	}
}
