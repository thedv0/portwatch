package snapshot

import (
	"testing"

	"github.com/wricardo/portwatch/internal/scanner"
)

func slPort(port int) scanner.Port {
	return scanner.Port{Port: port, Protocol: "tcp"}
}

func TestSlice_NoOptions_ReturnsAll(t *testing.T) {
	ports := []scanner.Port{slPort(80), slPort(443), slPort(8080)}
	result := Slice(ports, DefaultSliceOptions())
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
}

func TestSlice_Offset(t *testing.T) {
	ports := []scanner.Port{slPort(80), slPort(443), slPort(8080)}
	result := Slice(ports, SliceOptions{Offset: 1})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Port != 443 {
		t.Errorf("expected first port 443, got %d", result[0].Port)
	}
}

func TestSlice_Limit(t *testing.T) {
	ports := []scanner.Port{slPort(80), slPort(443), slPort(8080)}
	result := Slice(ports, SliceOptions{Limit: 2})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestSlice_OffsetBeyondLength(t *testing.T) {
	ports := []scanner.Port{slPort(80), slPort(443)}
	result := Slice(ports, SliceOptions{Offset: 10})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestSlice_Reverse(t *testing.T) {
	ports := []scanner.Port{slPort(80), slPort(443), slPort(8080)}
	result := Slice(ports, SliceOptions{Reverse: true})
	if result[0].Port != 8080 {
		t.Errorf("expected first port 8080, got %d", result[0].Port)
	}
}

func TestSlice_ReverseWithLimit(t *testing.T) {
	ports := []scanner.Port{slPort(80), slPort(443), slPort(8080)}
	result := Slice(ports, SliceOptions{Reverse: true, Limit: 2})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Port != 8080 {
		t.Errorf("expected first port 8080, got %d", result[0].Port)
	}
}

func TestSlice_Empty(t *testing.T) {
	result := Slice([]scanner.Port{}, DefaultSliceOptions())
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestSlice_DoesNotMutateInput(t *testing.T) {
	ports := []scanner.Port{slPort(80), slPort(443), slPort(8080)}
	Slice(ports, SliceOptions{Reverse: true})
	if ports[0].Port != 80 {
		t.Errorf("input slice was mutated")
	}
}
