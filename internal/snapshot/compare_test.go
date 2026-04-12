package snapshot

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func cport(proto string, port, pid int, proc string) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, PID: pid, Process: proc}
}

func TestCompare_NoChange(t *testing.T) {
	ports := []scanner.Port{cport("tcp", 80, 1, "nginx")}
	r := Compare(ports, ports)
	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Fatalf("expected no changes, got %+v", r)
	}
}

func TestCompare_Added(t *testing.T) {
	prev := []scanner.Port{cport("tcp", 80, 1, "nginx")}
	curr := []scanner.Port{cport("tcp", 80, 1, "nginx"), cport("tcp", 443, 2, "nginx")}
	r := Compare(prev, curr)
	if len(r.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(r.Added))
	}
	if r.Added[0].Port != 443 {
		t.Errorf("expected added port 443, got %d", r.Added[0].Port)
	}
}

func TestCompare_Removed(t *testing.T) {
	prev := []scanner.Port{cport("tcp", 80, 1, "nginx"), cport("udp", 53, 5, "dns")}
	curr := []scanner.Port{cport("tcp", 80, 1, "nginx")}
	r := Compare(prev, curr)
	if len(r.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(r.Removed))
	}
	if r.Removed[0].Port != 53 {
		t.Errorf("expected removed port 53, got %d", r.Removed[0].Port)
	}
}

func TestCompare_Changed(t *testing.T) {
	prev := []scanner.Port{cport("tcp", 8080, 10, "old-proc")}
	curr := []scanner.Port{cport("tcp", 8080, 99, "new-proc")}
	r := Compare(prev, curr)
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(r.Changed))
	}
	c := r.Changed[0]
	if c.OldPID != 10 || c.NewPID != 99 {
		t.Errorf("unexpected PID change: %d -> %d", c.OldPID, c.NewPID)
	}
	if c.OldProc != "old-proc" || c.NewProc != "new-proc" {
		t.Errorf("unexpected proc change: %s -> %s", c.OldProc, c.NewProc)
	}
}

func TestCompare_EmptyBoth(t *testing.T) {
	r := Compare(nil, nil)
	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Fatal("expected empty result for nil inputs")
	}
}

func TestCompareResult_Summary_NoChanges(t *testing.T) {
	r := CompareResult{}
	s := r.Summary()
	if s != "no changes detected" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestCompareResult_Summary_WithChanges(t *testing.T) {
	r := CompareResult{
		Added:   []scanner.Port{cport("tcp", 9090, 7, "app")},
		Removed: []scanner.Port{cport("udp", 161, 3, "snmpd")},
	}
	s := r.Summary()
	if !strings.Contains(s, "added 1") {
		t.Errorf("summary missing added count: %s", s)
	}
	if !strings.Contains(s, "removed 1") {
		t.Errorf("summary missing removed count: %s", s)
	}
	if !strings.Contains(s, "9090") {
		t.Errorf("summary missing added port: %s", s)
	}
}
