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

func makePivotResult() snapshot.PivotResult {
	return snapshot.PivotResult{
		GroupBy: "protocol",
		Groups: map[string][]snapshot.PortState{
			"tcp": {
				{Port: 80, Protocol: "tcp", Process: "nginx", PID: 100},
				{Port: 443, Protocol: "tcp", Process: "nginx", PID: 100},
			},
			"udp": {
				{Port: 53, Protocol: "udp", Process: "systemd-resolved", PID: 200},
			},
		},
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestBuildPivotReport_GroupCount(t *testing.T) {
	result := makePivotResult()
	r := report.BuildPivotReport(result)

	if r.GroupCount != 2 {
		t.Errorf("expected GroupCount=2, got %d", r.GroupCount)
	}
}

func TestBuildPivotReport_TotalPorts(t *testing.T) {
	result := makePivotResult()
	r := report.BuildPivotReport(result)

	if r.TotalPorts != 3 {
		t.Errorf("expected TotalPorts=3, got %d", r.TotalPorts)
	}
}

func TestBuildPivotReport_TimestampSet(t *testing.T) {
	result := makePivotResult()
	r := report.BuildPivotReport(result)

	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestBuildPivotReport_EmptyInput(t *testing.T) {
	result := snapshot.PivotResult{
		GroupBy: "protocol",
		Groups:  map[string][]snapshot.PortState{},
	}
	r := report.BuildPivotReport(result)

	if r.GroupCount != 0 {
		t.Errorf("expected GroupCount=0, got %d", r.GroupCount)
	}
	if r.TotalPorts != 0 {
		t.Errorf("expected TotalPorts=0, got %d", r.TotalPorts)
	}
}

func TestWritePivotText_ContainsGroupBy(t *testing.T) {
	result := makePivotResult()
	r := report.BuildPivotReport(result)

	var buf bytes.Buffer
	if err := report.WritePivotText(&buf, r); err != nil {
		t.Fatalf("WritePivotText error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "protocol") {
		t.Errorf("expected output to contain group-by field 'protocol', got:\n%s", out)
	}
}

func TestWritePivotText_ContainsGroupKeys(t *testing.T) {
	result := makePivotResult()
	r := report.BuildPivotReport(result)

	var buf bytes.Buffer
	if err := report.WritePivotText(&buf, r); err != nil {
		t.Fatalf("WritePivotText error: %v", err)
	}

	out := buf.String()
	for _, key := range []string{"tcp", "udp"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected output to contain group key %q", key)
		}
	}
}

func TestWritePivotJSON_ValidJSON(t *testing.T) {
	result := makePivotResult()
	r := report.BuildPivotReport(result)

	var buf bytes.Buffer
	if err := report.WritePivotJSON(&buf, r); err != nil {
		t.Fatalf("WritePivotJSON error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWritePivotJSON_GroupCountField(t *testing.T) {
	result := makePivotResult()
	r := report.BuildPivotReport(result)

	var buf bytes.Buffer
	if err := report.WritePivotJSON(&buf, r); err != nil {
		t.Fatalf("WritePivotJSON error: %v", err)
	}

	var out map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)

	val, ok := out["group_count"]
	if !ok {
		t.Fatal("expected 'group_count' field in JSON output")
	}
	if int(val.(float64)) != 2 {
		t.Errorf("expected group_count=2, got %v", val)
	}
}
