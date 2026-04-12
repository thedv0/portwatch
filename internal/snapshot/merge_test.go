package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func mp(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto}
}

func TestMerge_Union_NoDuplicates(t *testing.T) {
	left := []scanner.Port{mp(80, "tcp"), mp(443, "tcp")}
	right := []scanner.Port{mp(8080, "tcp"), mp(53, "udp")}
	got, err := Merge(left, right, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 4 {
		t.Errorf("expected 4 ports, got %d", len(got))
	}
}

func TestMerge_Union_DeduplicatesOverlap(t *testing.T) {
	left := []scanner.Port{mp(80, "tcp"), mp(443, "tcp")}
	right := []scanner.Port{mp(80, "tcp"), mp(8080, "tcp")}
	got, err := Merge(left, right, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 ports, got %d", len(got))
	}
}

func TestMerge_Intersect_CommonOnly(t *testing.T) {
	left := []scanner.Port{mp(80, "tcp"), mp(443, "tcp"), mp(22, "tcp")}
	right := []scanner.Port{mp(80, "tcp"), mp(22, "tcp"), mp(9090, "tcp")}
	got, err := Merge(left, right, MergeOptions{Strategy: MergeStrategyIntersect})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 ports, got %d", len(got))
	}
	for _, p := range got {
		if p.Port != 80 && p.Port != 22 {
			t.Errorf("unexpected port in intersect result: %d", p.Port)
		}
	}
}

func TestMerge_Intersect_EmptyResult(t *testing.T) {
	left := []scanner.Port{mp(80, "tcp")}
	right := []scanner.Port{mp(443, "tcp")}
	got, err := Merge(left, right, MergeOptions{Strategy: MergeStrategyIntersect})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 ports, got %d", len(got))
	}
}

func TestMerge_PreferLeft_OverridesRight(t *testing.T) {
	left := []scanner.Port{{Port: 80, Protocol: "tcp", Process: "nginx"}}
	right := []scanner.Port{{Port: 80, Protocol: "tcp", Process: "apache"}}
	got, err := Merge(left, right, MergeOptions{Strategy: MergeStrategyPreferLeft})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got))
	}
	if got[0].Process != "nginx" {
		t.Errorf("expected nginx, got %s", got[0].Process)
	}
}

func TestMerge_UnknownStrategy_ReturnsError(t *testing.T) {
	_, err := Merge(nil, nil, MergeOptions{Strategy: MergeStrategy(99)})
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestMerge_ResultIsSorted(t *testing.T) {
	left := []scanner.Port{mp(443, "tcp"), mp(80, "tcp")}
	right := []scanner.Port{mp(22, "tcp"), mp(8080, "tcp")}
	got, err := Merge(left, right, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(got); i++ {
		if got[i].Port < got[i-1].Port {
			t.Errorf("result not sorted at index %d: %d < %d", i, got[i].Port, got[i-1].Port)
		}
	}
}
