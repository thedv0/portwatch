package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func dp(proto string, port, pid int, process string) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, PID: pid, Process: process}
}

func TestDedupe_NoDuplicates(t *testing.T) {
	input := []scanner.Port{
		dp("tcp", 80, 100, "nginx"),
		dp("tcp", 443, 101, "nginx"),
		dp("udp", 53, 200, "dns"),
	}
	out := Dedupe(input, DefaultDedupeOptions())
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestDedupe_RemovesDuplicates(t *testing.T) {
	input := []scanner.Port{
		dp("tcp", 80, 100, "nginx"),
		dp("tcp", 80, 100, "nginx"),
		dp("tcp", 443, 101, "nginx"),
	}
	out := Dedupe(input, DefaultDedupeOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestDedupe_PreferHigherPID(t *testing.T) {
	input := []scanner.Port{
		dp("tcp", 80, 100, "nginx"),
		dp("tcp", 80, 999, "nginx"),
	}
	opts := DefaultDedupeOptions()
	opts.PreferHigherPID = true
	out := Dedupe(input, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].PID != 999 {
		t.Errorf("expected PID 999, got %d", out[0].PID)
	}
}

func TestDedupe_IgnoreProcess(t *testing.T) {
	input := []scanner.Port{
		dp("tcp", 80, 100, "nginx"),
		dp("tcp", 80, 101, "apache"),
	}
	opts := DefaultDedupeOptions()
	opts.IgnoreProcess = true
	out := Dedupe(input, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry (process ignored), got %d", len(out))
	}
}

func TestDedupe_DifferentProtocolNotDupe(t *testing.T) {
	input := []scanner.Port{
		dp("tcp", 53, 200, "dns"),
		dp("udp", 53, 200, "dns"),
	}
	out := Dedupe(input, DefaultDedupeOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries (different protocols), got %d", len(out))
	}
}

func TestDedupe_Empty(t *testing.T) {
	out := Dedupe([]scanner.Port{}, DefaultDedupeOptions())
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}

func TestDefaultDedupeOptions_Values(t *testing.T) {
	opts := DefaultDedupeOptions()
	if opts.PreferHigherPID {
		t.Error("expected PreferHigherPID to be false")
	}
	if opts.IgnoreProcess {
		t.Error("expected IgnoreProcess to be false")
	}
}
