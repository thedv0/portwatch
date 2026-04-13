package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func fport(proto string, port int) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port}
}

func TestFlatten_CombinesAllSnaps(t *testing.T) {
	snaps := [][]scanner.Port{
		{fport("tcp", 80), fport("tcp", 443)},
		{fport("udp", 53)},
	}
	opts := DefaultFlattenOptions()
	opts.Deduplicate = false
	result := Flatten(snaps, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(result))
	}
}

func TestFlatten_DeduplicatesOverlap(t *testing.T) {
	snaps := [][]scanner.Port{
		{fport("tcp", 80), fport("tcp", 443)},
		{fport("tcp", 80), fport("udp", 53)},
	}
	opts := DefaultFlattenOptions()
	result := Flatten(snaps, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3 unique ports, got %d", len(result))
	}
}

func TestFlatten_SortByPort(t *testing.T) {
	snaps := [][]scanner.Port{
		{fport("tcp", 8080), fport("tcp", 22), fport("tcp", 443)},
	}
	opts := DefaultFlattenOptions()
	result := Flatten(snaps, opts)
	if result[0].Port != 22 || result[1].Port != 443 || result[2].Port != 8080 {
		t.Fatalf("expected sorted order 22,443,8080 got %v", result)
	}
}

func TestFlatten_FilterByProtocol(t *testing.T) {
	snaps := [][]scanner.Port{
		{fport("tcp", 80), fport("udp", 53), fport("tcp", 443)},
	}
	opts := DefaultFlattenOptions()
	opts.IncludeProtocols = []string{"tcp"}
	result := Flatten(snaps, opts)
	for _, p := range result {
		if p.Protocol != "tcp" {
			t.Errorf("expected only tcp, got %s", p.Protocol)
		}
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 tcp ports, got %d", len(result))
	}
}

func TestFlatten_EmptyInput(t *testing.T) {
	result := Flatten(nil, DefaultFlattenOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestFlatten_NoSort(t *testing.T) {
	snaps := [][]scanner.Port{
		{fport("tcp", 9000), fport("tcp", 22)},
	}
	opts := DefaultFlattenOptions()
	opts.SortByPort = false
	result := Flatten(snaps, opts)
	if result[0].Port != 9000 {
		t.Errorf("expected original order preserved, first port should be 9000")
	}
}
