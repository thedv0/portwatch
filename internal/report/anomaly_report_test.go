package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/snapshot"
)

var fixedAnomalyReportTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeAnomalies() []snapshot.Anomaly {
	return []snapshot.Anomaly{
		{Type: snapshot.AnomalyNewPort, Port: 9090, Protocol: "tcp", Process: "app", PID: 42, Message: "new port opened", DetectedAt: fixedAnomalyReportTime},
		{Type: snapshot.AnomalyGonePort, Port: 8080, Protocol: "tcp", Process: "old", PID: 10, Message: "port closed", DetectedAt: fixedAnomalyReportTime},
		{Type: snapshot.AnomalyPIDChanged, Port: 443, Protocol: "tcp", Process: "nginx", PID: 99, Message: "PID changed on port", DetectedAt: fixedAnomalyReportTime},
	}
}

func TestBuildAnomalyReport_Counts(t *testing.T) {
	r := BuildAnomalyReport(makeAnomalies(), fixedAnomalyReportTime)
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
	if r.NewPorts != 1 {
		t.Errorf("expected NewPorts=1, got %d", r.NewPorts)
	}
	if r.GonePorts != 1 {
		t.Errorf("expected GonePorts=1, got %d", r.GonePorts)
	}
	if r.PIDChanges != 1 {
		t.Errorf("expected PIDChanges=1, got %d", r.PIDChanges)
	}
}

func TestBuildAnomalyReport_TimestampSet(t *testing.T) {
	r := BuildAnomalyReport(nil, fixedAnomalyReportTime)
	if !r.Timestamp.Equal(fixedAnomalyReportTime) {
		t.Errorf("timestamp mismatch: %v", r.Timestamp)
	}
}

func TestBuildAnomalyReport_EmptyInput(t *testing.T) {
	r := BuildAnomalyReport(nil, fixedAnomalyReportTime)
	if r.Total != 0 {
		t.Errorf("expected Total=0, got %d", r.Total)
	}
}

func TestWriteAnomalyText_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	r := BuildAnomalyReport(makeAnomalies(), fixedAnomalyReportTime)
	if err := WriteAnomalyText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Anomaly Report", "New:", "Gone:", "9090", "8080"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestWriteAnomalyText_NoAnomalies(t *testing.T) {
	var buf bytes.Buffer
	r := BuildAnomalyReport(nil, fixedAnomalyReportTime)
	WriteAnomalyText(&buf, r)
	if !strings.Contains(buf.String(), "No anomalies") {
		t.Error("expected 'No anomalies' message")
	}
}

func TestWriteAnomalyJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := BuildAnomalyReport(makeAnomalies(), fixedAnomalyReportTime)
	if err := WriteAnomalyJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out AnomalyReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Total != 3 {
		t.Errorf("expected Total=3 in JSON, got %d", out.Total)
	}
}
