package snapshot

import (
	"testing"
)

func pport(num int, proto, process string) PortState {
	return PortState{Port: num, Protocol: proto, Process: process, PID: 100}
}

func TestApplyPatches_SetProcess(t *testing.T) {
	ports := []PortState{pport(80, "tcp", "nginx")}
	patches := []Patch{{Op: PatchSet, Key: "process", Value: "apache", PortNum: 80, Proto: "tcp"}}

	res, err := ApplyPatches(ports, patches, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Applied != 1 {
		t.Errorf("expected 1 applied, got %d", res.Applied)
	}
	if res.Ports[0].Process != "apache" {
		t.Errorf("expected process=apache, got %s", res.Ports[0].Process)
	}
}

func TestApplyPatches_DeleteProcess(t *testing.T) {
	ports := []PortState{pport(443, "tcp", "nginx")}
	patches := []Patch{{Op: PatchDelete, Key: "process", PortNum: 443, Proto: "tcp"}}

	res, err := ApplyPatches(ports, patches, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Ports[0].Process != "" {
		t.Errorf("expected empty process, got %s", res.Ports[0].Process)
	}
}

func TestApplyPatches_MissingPort_IgnoreMissing(t *testing.T) {
	ports := []PortState{pport(80, "tcp", "nginx")}
	patches := []Patch{{Op: PatchSet, Key: "process", Value: "x", PortNum: 9999, Proto: "tcp"}}

	res, err := ApplyPatches(ports, patches, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestApplyPatches_MissingPort_ErrorOnMissing(t *testing.T) {
	ports := []PortState{pport(80, "tcp", "nginx")}
	patches := []Patch{{Op: PatchSet, Key: "process", Value: "x", PortNum: 9999, Proto: "tcp"}}
	opts := PatchOptions{IgnoreMissing: false}

	_, err := ApplyPatches(ports, patches, opts)
	if err == nil {
		t.Error("expected error for missing port, got nil")
	}
}

func TestApplyPatches_UnknownKey_Skipped(t *testing.T) {
	ports := []PortState{pport(22, "tcp", "sshd")}
	patches := []Patch{{Op: PatchSet, Key: "unknown_field", Value: "x", PortNum: 22, Proto: "tcp"}}

	res, err := ApplyPatches(ports, patches, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestApplyPatches_TimestampSet(t *testing.T) {
	ports := []PortState{pport(80, "tcp", "nginx")}
	res, _ := ApplyPatches(ports, nil, DefaultPatchOptions())
	if res.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestApplyPatches_DoesNotMutateOriginal(t *testing.T) {
	original := []PortState{pport(80, "tcp", "nginx")}
	patches := []Patch{{Op: PatchSet, Key: "process", Value: "changed", PortNum: 80, Proto: "tcp"}}

	_, _ = ApplyPatches(original, patches, DefaultPatchOptions())
	if original[0].Process != "nginx" {
		t.Error("original slice was mutated")
	}
}
