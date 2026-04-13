package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func pvPort(port int, proto, process string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: process}
}

func makeSnapsForPivot() []Snapshot {
	return []Snapshot{
		{Ports: []scanner.Port{
			pvPort(80, "tcp", "nginx"),
			pvPort(443, "tcp", "nginx"),
			pvPort(53, "udp", "dnsmasq"),
		}},
		{Ports: []scanner.Port{
			pvPort(8080, "tcp", "caddy"),
			pvPort(53, "udp", "dnsmasq"),
		}},
	}
}

func TestPivot_ByProtocol_GroupsCorrectly(t *testing.T) {
	snaps := makeSnapsForPivot()
	opts := DefaultPivotOptions()
	res := Pivot(snaps, opts)

	if res.Field != PivotByProtocol {
		t.Fatalf("expected field %q, got %q", PivotByProtocol, res.Field)
	}
	if len(res.Buckets["tcp"]) != 3 {
		t.Errorf("expected 3 tcp ports, got %d", len(res.Buckets["tcp"]))
	}
	if len(res.Buckets["udp"]) != 2 {
		t.Errorf("expected 2 udp ports, got %d", len(res.Buckets["udp"]))
	}
}

func TestPivot_ByProcess_GroupsCorrectly(t *testing.T) {
	snaps := makeSnapsForPivot()
	opts := PivotOptions{Field: PivotByProcess, SortKeys: true}
	res := Pivot(snaps, opts)

	if len(res.Buckets["nginx"]) != 2 {
		t.Errorf("expected 2 nginx ports, got %d", len(res.Buckets["nginx"]))
	}
	if len(res.Buckets["dnsmasq"]) != 2 {
		t.Errorf("expected 2 dnsmasq ports, got %d", len(res.Buckets["dnsmasq"]))
	}
	if len(res.Buckets["caddy"]) != 1 {
		t.Errorf("expected 1 caddy port, got %d", len(res.Buckets["caddy"]))
	}
}

func TestPivot_ByPort_KeyIsPortNumber(t *testing.T) {
	snaps := []Snapshot{
		{Ports: []scanner.Port{pvPort(80, "tcp", "nginx")}},
	}
	res := Pivot(snaps, PivotOptions{Field: PivotByPort, SortKeys: false})

	if _, ok := res.Buckets["80"]; !ok {
		t.Error("expected bucket for port 80")
	}
}

func TestPivot_SortKeys_Ordered(t *testing.T) {
	snaps := makeSnapsForPivot()
	res := Pivot(snaps, PivotOptions{Field: PivotByProtocol, SortKeys: true})

	if len(res.Keys) < 2 {
		t.Fatal("expected at least 2 keys")
	}
	for i := 1; i < len(res.Keys); i++ {
		if res.Keys[i] < res.Keys[i-1] {
			t.Errorf("keys not sorted: %v", res.Keys)
		}
	}
}

func TestPivot_EmptySnaps_ReturnsEmpty(t *testing.T) {
	res := Pivot(nil, DefaultPivotOptions())
	if len(res.Buckets) != 0 {
		t.Errorf("expected empty buckets, got %d", len(res.Buckets))
	}
	if len(res.Keys) != 0 {
		t.Errorf("expected empty keys, got %d", len(res.Keys))
	}
}

func TestPivot_UnknownProcess_FallbackKey(t *testing.T) {
	snaps := []Snapshot{
		{Ports: []scanner.Port{pvPort(9999, "tcp", "")}},
	}
	res := Pivot(snaps, PivotOptions{Field: PivotByProcess, SortKeys: false})
	if _, ok := res.Buckets["(unknown)"]; !ok {
		t.Error("expected fallback key '(unknown)' for empty process")
	}
}
