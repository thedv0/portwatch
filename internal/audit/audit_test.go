package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	err := l.Log(LevelInfo, "daemon_started", map[string]string{"version": "1.0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if e.Event != "daemon_started" {
		t.Errorf("expected event daemon_started, got %s", e.Event)
	}
	if e.Level != LevelInfo {
		t.Errorf("expected level INFO, got %s", e.Level)
	}
	if e.Details["version"] != "1.0" {
		t.Errorf("expected version 1.0 in details")
	}
}

func TestLog_TimestampSet(t *testing.T) {
	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	var buf bytes.Buffer
	l := NewLogger(&buf)
	l.clock = func() time.Time { return fixed }
	_ = l.Info("test", nil)
	var e Entry
	_ = json.Unmarshal(buf.Bytes(), &e)
	if !e.Timestamp.Equal(fixed) {
		t.Errorf("expected timestamp %v, got %v", fixed, e.Timestamp)
	}
}

func TestWarn_SetsLevel(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	_ = l.Warn("unexpected_port", nil)
	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("expected WARN in output: %s", buf.String())
	}
}

func TestAlert_SetsLevel(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	_ = l.Alert("port_scan_detected", map[string]string{"port": "4444"})
	if !strings.Contains(buf.String(), "ALERT") {
		t.Errorf("expected ALERT in output: %s", buf.String())
	}
}

func TestNewLogger_NilUsesStderr(t *testing.T) {
	l := NewLogger(nil)
	if l.w == nil {
		t.Error("expected non-nil writer")
	}
}

func TestLog_NilDetails(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	err := l.Info("no_details", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(e.Details) != 0 {
		t.Errorf("expected empty details, got %v", e.Details)
	}
}
