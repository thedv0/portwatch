package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/snapshot"
)

func makeReport() report.Report {
	return report.Report{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Added: []snapshot.Port{
			{Port: 8080, Protocol: "tcp", PID: 1234},
		},
		Removed: []snapshot.Port{
			{Port: 9090, Protocol: "tcp", PID: 5678},
		},
		Total: 5,
	}
}

func TestWriteText_ContainsExpectedLines(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	if err := w.Write(makeReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "total open: 5") {
		t.Errorf("expected total open in output, got: %s", out)
	}
	if !strings.Contains(out, "+ tcp/8080") {
		t.Errorf("expected added port in output, got: %s", out)
	}
	if !strings.Contains(out, "- tcp/9090") {
		t.Errorf("expected removed port in output, got: %s", out)
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)
	if err := w.Write(makeReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var r report.Report
	if err := json.Unmarshal(buf.Bytes(), &r); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if r.Total != 5 {
		t.Errorf("expected Total=5, got %d", r.Total)
	}
	if len(r.Added) != 1 || r.Added[0].Port != 8080 {
		t.Errorf("unexpected Added ports: %+v", r.Added)
	}
}

func TestWriteJSON_TimestampPreserved(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)
	if err := w.Write(makeReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var r report.Report
	if err := json.Unmarshal(buf.Bytes(), &r); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	want := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !r.Timestamp.Equal(want) {
		t.Errorf("expected Timestamp=%v, got %v", want, r.Timestamp)
	}
}

func TestNewWriter_NilUsesStdout(t *testing.T) {
	// Should not panic when out is nil
	w := report.NewWriter(nil, report.FormatText)
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestNewWriter_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, "")
	if err := w.Write(makeReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Default format is text, should not start with '{'
	if strings.HasPrefix(buf.String(), "{") {
		t.Error("expected text format, got JSON-like output")
	}
}
