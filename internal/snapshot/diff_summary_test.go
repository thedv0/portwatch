package snapshot

import (
	"testing"
	"time"
)

func makeDiffPort(proto string, port int) Port {
	return Port{Protocol: proto, Port: port, PID: 100, Process: "test"}
}

func fixedDiffClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestSummarizeDiff_NoChanges(t *testing.T) {
	d := DiffResult{
		Unchanged: []Port{makeDiffPort("tcp", 80)},
	}
	s := SummarizeDiff(d, nil)
	if s.HasChanges {
		t.Error("expected no changes")
	}
	if s.UnchangedCount != 1 {
		t.Errorf("expected 1 unchanged, got %d", s.UnchangedCount)
	}
}

func TestSummarizeDiff_Added(t *testing.T) {
	d := DiffResult{
		Added: []Port{makeDiffPort("tcp", 8080), makeDiffPort("udp", 53)},
	}
	s := SummarizeDiff(d, nil)
	if !s.HasChanges {
		t.Error("expected changes")
	}
	if s.AddedCount != 2 {
		t.Errorf("expected 2 added, got %d", s.AddedCount)
	}
	if len(s.AddedPorts) != 2 {
		t.Errorf("expected 2 added port keys, got %d", len(s.AddedPorts))
	}
}

func TestSummarizeDiff_Removed(t *testing.T) {
	d := DiffResult{
		Removed: []Port{makeDiffPort("tcp", 443)},
	}
	s := SummarizeDiff(d, nil)
	if !s.HasChanges {
		t.Error("expected changes")
	}
	if s.RemovedCount != 1 {
		t.Errorf("expected 1 removed, got %d", s.RemovedCount)
	}
}

func TestSummarizeDiff_TimestampSet(t *testing.T) {
	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	d := DiffResult{}
	s := SummarizeDiff(d, fixedDiffClock(fixed))
	if !s.Timestamp.Equal(fixed) {
		t.Errorf("expected timestamp %v, got %v", fixed, s.Timestamp)
	}
}

func TestSummarizeDiff_PortKeys(t *testing.T) {
	d := DiffResult{
		Added:   []Port{makeDiffPort("tcp", 22)},
		Removed: []Port{makeDiffPort("udp", 123)},
	}
	s := SummarizeDiff(d, nil)
	if len(s.AddedPorts) == 0 || s.AddedPorts[0] != "tcp:22" {
		t.Errorf("unexpected added port key: %v", s.AddedPorts)
	}
	if len(s.RemovedPorts) == 0 || s.RemovedPorts[0] != "udp:123" {
		t.Errorf("unexpected removed port key: %v", s.RemovedPorts)
	}
}
