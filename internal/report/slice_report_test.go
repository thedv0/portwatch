package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wricardo/portwatch/internal/scanner"
	"github.com/wricardo/portwatch/internal/snapshot"
)

func makeSlicePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Protocol: "tcp", PID: 100, Process: "nginx"},
		{Port: 443, Protocol: "tcp", PID: 101, Process: "nginx"},
		{Port: 8080, Protocol: "tcp", PID: 200, Process: "app"},
		{Port: 9090, Protocol: "udp", PID: 300, Process: "prom"},
	}
}

func TestBuildSliceReport_Counts(t *testing.T) {
	ports := makeSlicePorts()
	opts := snapshot.SliceOptions{Offset: 1, Limit: 2}
	r := BuildSliceReport(ports, opts)
	if r.Total != 4 {
		t.Errorf("expected total 4, got %d", r.Total)
	}
	if r.Returned != 2 {
		t.Errorf("expected returned 2, got %d", r.Returned)
	}
}

func TestBuildSliceReport_TimestampSet(t *testing.T) {
	r := BuildSliceReport(makeSlicePorts(), snapshot.DefaultSliceOptions())
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestWriteSliceText_ContainsHeaders(t *testing.T) {
	r := BuildSliceReport(makeSlicePorts(), snapshot.SliceOptions{Limit: 2})
	var buf bytes.Buffer
	if err := WriteSliceText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Slice Report") {
		t.Error("expected 'Slice Report' header")
	}
	if !strings.Contains(out, "Input:") {
		t.Error("expected 'Input:' line")
	}
	if !strings.Contains(out, "Returned:") {
		t.Error("expected 'Returned:' line")
	}
}

func TestWriteSliceText_ContainsPortNumbers(t *testing.T) {
	ports := makeSlicePorts()
	r := BuildSliceReport(ports, snapshot.DefaultSliceOptions())
	var buf bytes.Buffer
	WriteSliceText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
}

func TestWriteSliceJSON_ValidJSON(t *testing.T) {
	r := BuildSliceReport(makeSlicePorts(), snapshot.DefaultSliceOptions())
	var buf bytes.Buffer
	if err := WriteSliceJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out SliceReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Total != 4 {
		t.Errorf("expected total 4, got %d", out.Total)
	}
}

func TestBuildSliceReport_EmptyInput(t *testing.T) {
	r := BuildSliceReport([]scanner.Port{}, snapshot.DefaultSliceOptions())
	if r.Total != 0 || r.Returned != 0 {
		t.Errorf("expected zero counts, got total=%d returned=%d", r.Total, r.Returned)
	}
}
