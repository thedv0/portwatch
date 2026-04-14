package snapshot

import (
	"testing"
)

func maskPort(port int, pid int, process, proto string) PortEntry {
	return PortEntry{Port: port, PID: pid, Process: process, Protocol: proto}
}

func TestMask_DefaultOptions_MasksProcess(t *testing.T) {
	ports := []PortEntry{maskPort(80, 100, "nginx", "tcp")}
	out := Mask(ports, DefaultMaskOptions())
	if out[0].Process != "[redacted]" {
		t.Errorf("expected process redacted, got %q", out[0].Process)
	}
	if out[0].PID != 100 {
		t.Errorf("expected PID preserved, got %d", out[0].PID)
	}
	if out[0].Port != 80 {
		t.Errorf("expected port preserved, got %d", out[0].Port)
	}
}

func TestMask_MaskPID(t *testing.T) {
	ports := []PortEntry{maskPort(443, 999, "sshd", "tcp")}
	opts := DefaultMaskOptions()
	opts.MaskPID = true
	out := Mask(ports, opts)
	if out[0].PID != 0 {
		t.Errorf("expected PID zeroed, got %d", out[0].PID)
	}
}

func TestMask_MaskPort(t *testing.T) {
	ports := []PortEntry{maskPort(8080, 42, "app", "tcp")}
	opts := DefaultMaskOptions()
	opts.MaskPort = true
	out := Mask(ports, opts)
	if out[0].Port != 0 {
		t.Errorf("expected port zeroed, got %d", out[0].Port)
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	ports := []PortEntry{maskPort(22, 5, "sshd", "tcp")}
	opts := DefaultMaskOptions()
	opts.Placeholder = "***"
	out := Mask(ports, opts)
	if out[0].Process != "***" {
		t.Errorf("expected placeholder ***, got %q", out[0].Process)
	}
}

func TestMask_EmptyPlaceholderFallsBack(t *testing.T) {
	ports := []PortEntry{maskPort(22, 5, "sshd", "tcp")}
	opts := DefaultMaskOptions()
	opts.Placeholder = "   "
	out := Mask(ports, opts)
	if out[0].Process != "[redacted]" {
		t.Errorf("expected fallback placeholder, got %q", out[0].Process)
	}
}

func TestMask_OriginalUnchanged(t *testing.T) {
	original := []PortEntry{maskPort(80, 10, "nginx", "tcp")}
	Mask(original, DefaultMaskOptions())
	if original[0].Process != "nginx" {
		t.Error("original slice should not be modified")
	}
}

func TestMask_EmptyInput(t *testing.T) {
	out := Mask([]PortEntry{}, DefaultMaskOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
