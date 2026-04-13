package snapshot

import (
	"testing"
)

func cport(port int, proto, process string, pid int) Port {
	return Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func TestCorrelate_ByPort_GroupsMatchingPorts(t *testing.T) {
	snap1 := []Port{cport(80, "tcp", "nginx", 100)}
	snap2 := []Port{cport(80, "tcp", "nginx", 200)}

	opts := DefaultCorrelateOptions()
	groups := Correlate([][]Port{snap1, snap2}, opts)

	var found *CorrelatedGroup
	for i := range groups {
		if groups[i].Key == "port:tcp:80" {
			found = &groups[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected group for port:tcp:80")
	}
	if found.Count != 2 {
		t.Errorf("expected count 2, got %d", found.Count)
	}
}

func TestCorrelate_ByProcess_GroupsSameProcess(t *testing.T) {
	snap1 := []Port{cport(443, "tcp", "nginx", 101)}
	snap2 := []Port{cport(8443, "tcp", "nginx", 102)}

	opts := CorrelateOptions{MatchByProcess: true}
	groups := Correlate([][]Port{snap1, snap2}, opts)

	var found *CorrelatedGroup
	for i := range groups {
		if groups[i].Key == "process:nginx" {
			found = &groups[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected group for process:nginx")
	}
	if found.Count != 2 {
		t.Errorf("expected count 2, got %d", found.Count)
	}
}

func TestCorrelate_ByPID_GroupsSamePID(t *testing.T) {
	snap1 := []Port{cport(80, "tcp", "app", 999)}
	snap2 := []Port{cport(81, "tcp", "app", 999)}

	opts := CorrelateOptions{MatchByPID: true}
	groups := Correlate([][]Port{snap1, snap2}, opts)

	var found *CorrelatedGroup
	for i := range groups {
		if groups[i].Key == "pid:999" {
			found = &groups[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected group for pid:999")
	}
	if found.Count != 2 {
		t.Errorf("expected count 2, got %d", found.Count)
	}
}

func TestCorrelate_EmptySnaps_ReturnsEmpty(t *testing.T) {
	groups := Correlate([][]Port{}, DefaultCorrelateOptions())
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestCorrelate_NoOptions_ReturnsEmpty(t *testing.T) {
	snap := []Port{cport(80, "tcp", "nginx", 1)}
	groups := Correlate([][]Port{snap}, CorrelateOptions{})
	if len(groups) != 0 {
		t.Errorf("expected no groups with all options disabled, got %d", len(groups))
	}
}

func TestCorrelate_SortedByCountDescending(t *testing.T) {
	snap1 := []Port{cport(80, "tcp", "a", 1), cport(443, "tcp", "b", 2)}
	snap2 := []Port{cport(80, "tcp", "a", 3)}

	opts := CorrelateOptions{MatchByPort: true}
	groups := Correlate([][]Port{snap1, snap2}, opts)

	if len(groups) < 2 {
		t.Fatal("expected at least 2 groups")
	}
	if groups[0].Count < groups[1].Count {
		t.Errorf("groups not sorted by count descending: %d < %d", groups[0].Count, groups[1].Count)
	}
}
