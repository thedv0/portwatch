package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeMaskedPorts() []snapshot.PortEntry {
	return []snapshot.PortEntry{
		{Port: 80, PID: 0, Protocol: "tcp", Process: "[redacted]"},
		{Port: 443, PID: 0, Protocol: "tcp", Process: "[redacted]"},
	}
}

func TestBuildMaskReport_Total(t *testing.T) {
	ports := makeMaskedPorts()
	r := BuildMaskReport(ports, snapshot.DefaultMaskOptions())
	if r.Total != 2 {
		t.Errorf("expected total 2, got %d", r.Total)
	}
}

func TestBuildMaskReport_TimestampSet(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	r := BuildMaskReport(makeMaskedPorts(), snapshot.DefaultMaskOptions())
	if r.Timestamp.Before(before) {
		t.Error("timestamp should be recent")
	}
}

func TestBuildMaskReport_OptionsPreserved(t *testing.T) {
	opts := snapshot.DefaultMaskOptions()
	opts.MaskPID = true
	r := BuildMaskReport(makeMaskedPorts(), opts)
	if !r.Options.MaskPID {
		t.Error("expected MaskPID to be true in report")
	}
}

func TestWriteMaskText_ContainsHeaders(t *testing.T) {
	r := BuildMaskReport(makeMaskedPorts(), snapshot.DefaultMaskOptions())
	var buf bytes.Buffer
	if err := WriteMaskText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Mask Report", "Total entries", "MaskProcess", "---"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteMaskText_ContainsPortNumbers(t *testing.T) {
	r := BuildMaskReport(makeMaskedPorts(), snapshot.DefaultMaskOptions())
	var buf bytes.Buffer
	WriteMaskText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "80") || !strings.Contains(out, "443") {
		t.Error("expected port numbers in text output")
	}
}

func TestWriteMaskJSON_ValidJSON(t *testing.T) {
	r := BuildMaskReport(makeMaskedPorts(), snapshot.DefaultMaskOptions())
	var buf bytes.Buffer
	if err := WriteMaskJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out MaskReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Total != 2 {
		t.Errorf("expected total 2 in JSON, got %d", out.Total)
	}
}
