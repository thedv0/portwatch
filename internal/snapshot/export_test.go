package snapshot

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSnap() Snapshot {
	return Snapshot{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Ports: []scanner.PortState{
			{Protocol: "tcp", Port: 80, PID: 100, Process: "nginx"},
			{Protocol: "udp", Port: 53, PID: 200, Process: "systemd-resolved"},
		},
	}
}

func TestExport_JSON_Valid(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, FormatJSON)
	if err := ex.Write(makeSnap()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
	if records[0].Port.Process != "nginx" {
		t.Errorf("expected nginx, got %s", records[0].Port.Process)
	}
}

func TestExport_CSV_HeaderAndRows(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, FormatCSV)
	if err := ex.Write(makeSnap()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header+2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp,protocol") {
		t.Errorf("unexpected header: %s", lines[0])
	}
	if !strings.Contains(lines[1], "nginx") {
		t.Errorf("expected nginx in row 1: %s", lines[1])
	}
}

func TestExport_DefaultFormat_IsJSON(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, "")
	if ex.format != FormatJSON {
		t.Errorf("expected default format JSON, got %q", ex.format)
	}
}

func TestExport_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(&buf, "xml")
	if err := ex.Write(makeSnap()); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestExport_EmptyPorts(t *testing.T) {
	snap := Snapshot{Timestamp: time.Now(), Ports: []scanner.PortState{}}
	var buf bytes.Buffer
	ex := NewExporter(&buf, FormatJSON)
	if err := ex.Write(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}
