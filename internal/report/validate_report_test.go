package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeValidatePorts() []snapshot.PortEntry {
	return []snapshot.PortEntry{
		{Port: 80, PID: 100, Process: "nginx", Protocol: "tcp"},
		{Port: 8080, PID: 0, Process: "", Protocol: "tcp"},
		{Port: 99999, PID: 1, Process: "bad", Protocol: "tcp"},
	}
}

func TestBuildValidateReport_Counts(t *testing.T) {
	ports := makeValidatePorts()
	opts := snapshot.DefaultValidateOptions()
	opts.RequireProcess = true
	r, err := BuildValidateReport(ports, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
	if r.ErrorCount == 0 {
		t.Error("expected at least one error")
	}
	if r.WarningCount == 0 {
		t.Error("expected at least one warning")
	}
}

func TestBuildValidateReport_TimestampSet(t *testing.T) {
	ports := []snapshot.PortEntry{{Port: 80, PID: 1, Process: "svc", Protocol: "tcp"}}
	r, err := BuildValidateReport(ports, snapshot.DefaultValidateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestWriteValidateText_ContainsHeaders(t *testing.T) {
	ports := makeValidatePorts()
	r, _ := BuildValidateReport(ports, snapshot.DefaultValidateOptions())
	var buf bytes.Buffer
	WriteValidateText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Validation Report") {
		t.Error("expected 'Validation Report' header")
	}
	if !strings.Contains(out, "Total ports") {
		t.Error("expected 'Total ports' line")
	}
}

func TestWriteValidateText_NoIssues(t *testing.T) {
	ports := []snapshot.PortEntry{{Port: 80, PID: 1, Process: "nginx", Protocol: "tcp"}}
	r, _ := BuildValidateReport(ports, snapshot.DefaultValidateOptions())
	var buf bytes.Buffer
	WriteValidateText(&buf, r)
	if !strings.Contains(buf.String(), "No issues found") {
		t.Error("expected 'No issues found' message")
	}
}

func TestWriteValidateJSON_ValidJSON(t *testing.T) {
	ports := makeValidatePorts()
	r, _ := BuildValidateReport(ports, snapshot.DefaultValidateOptions())
	var buf bytes.Buffer
	if err := WriteValidateJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["timestamp"]; !ok {
		t.Error("expected 'timestamp' field in JSON")
	}
	if _, ok := out["issues"]; !ok {
		t.Error("expected 'issues' field in JSON")
	}
}

func TestBuildValidateReport_EmptyInput(t *testing.T) {
	r, err := BuildValidateReport(nil, snapshot.DefaultValidateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Total != 0 {
		t.Errorf("expected Total=0, got %d", r.Total)
	}
	if len(r.Issues) != 0 {
		t.Errorf("expected no issues for empty input")
	}
}
