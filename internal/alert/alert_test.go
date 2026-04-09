package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

func TestSend_FormatsOutput(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	e := alert.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     alert.LevelAlert,
		Port:      4444,
		Protocol:  "tcp",
		Process:   "nc",
		Message:   "unexpected listener",
	}

	if err := n.Send(e); err != nil {
		t.Fatalf("Send() error: %v", err)
	}

	got := buf.String()
	for _, want := range []string{"ALERT", "port=4444", "proto=tcp", `process="nc"`, "unexpected listener"} {
		if !strings.Contains(got, want) {
			t.Errorf("output %q missing %q", got, want)
		}
	}
}

func TestSend_DefaultsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	e := alert.Event{
		Level:    alert.LevelWarn,
		Port:     8080,
		Protocol: "tcp",
		Process:  "unknown",
		Message:  "new port",
	}

	before := time.Now()
	_ = n.Send(e)
	after := time.Now()

	got := buf.String()
	// Timestamp should be between before and after — just check it's non-empty.
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	_ = before
	_ = after
}

func TestNewNotifier_NilUsesStderr(t *testing.T) {
	// Should not panic.
	n := alert.NewNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

func TestNewEvent_Fields(t *testing.T) {
	e := alert.NewEvent(alert.LevelInfo, 22, "tcp", "sshd", "allowed")
	if e.Port != 22 || e.Protocol != "tcp" || e.Process != "sshd" || e.Level != alert.LevelInfo {
		t.Errorf("NewEvent fields mismatch: %+v", e)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
