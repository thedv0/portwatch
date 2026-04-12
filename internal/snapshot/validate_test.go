package snapshot

import (
	"testing"
)

func vport(port, pid int, process string) PortEntry {
	return PortEntry{Port: port, PID: pid, Process: process, Protocol: "tcp"}
}

func TestValidate_NoIssues(t *testing.T) {
	ports := []PortEntry{vport(80, 1234, "nginx"), vport(443, 5678, "nginx")}
	opts := DefaultValidateOptions()
	res, err := Validate(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(res.Issues))
	}
}

func TestValidate_PortOutOfRange(t *testing.T) {
	ports := []PortEntry{vport(99999, 100, "bad")}
	opts := DefaultValidateOptions()
	res, err := Validate(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasErrors() {
		t.Error("expected error-level issue for out-of-range port")
	}
}

func TestValidate_PIDZeroWarning(t *testing.T) {
	ports := []PortEntry{vport(8080, 0, "unknown")}
	opts := DefaultValidateOptions()
	res, err := Validate(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(res.Issues))
	}
	if res.Issues[0].Level != LevelWarning {
		t.Errorf("expected warning, got %d", res.Issues[0].Level)
	}
}

func TestValidate_PIDZeroAllowed(t *testing.T) {
	ports := []PortEntry{vport(8080, 0, "unknown")}
	opts := DefaultValidateOptions()
	opts.AllowPIDZero = true
	res, err := Validate(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Issues) != 0 {
		t.Errorf("expected no issues when PIDZero allowed, got %d", len(res.Issues))
	}
}

func TestValidate_RequireProcess(t *testing.T) {
	ports := []PortEntry{vport(3000, 42, "")}
	opts := DefaultValidateOptions()
	opts.RequireProcess = true
	res, err := Validate(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Issues) == 0 {
		t.Error("expected warning for missing process name")
	}
	if res.Issues[0].Level != LevelWarning {
		t.Errorf("expected warning level, got %d", res.Issues[0].Level)
	}
}

func TestValidate_InvalidMaxPort(t *testing.T) {
	ports := []PortEntry{vport(80, 1, "svc")}
	opts := DefaultValidateOptions()
	opts.MaxPort = 0
	_, err := Validate(ports, opts)
	if err == nil {
		t.Error("expected error for MaxPort=0")
	}
}

func TestValidate_HasErrors_False(t *testing.T) {
	res := &ValidationResult{Issues: []ValidationIssue{{Level: LevelWarning}}}
	if res.HasErrors() {
		t.Error("expected HasErrors to be false with only warnings")
	}
}
