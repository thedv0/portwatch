package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeRollupResult() snapshot.RollupResult {
	return snapshot.RollupResult{
		Timestamp:     time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		SnapshotCount: 5,
		Dropped:       2,
		Ports: []scanner.Port{
			{Port: 80, Protocol: "tcp", Process: "nginx", PID: 1001},
			{Port: 443, Protocol: "tcp", Process: "nginx", PID: 1001},
			{Port: 22, Protocol: "tcp", Process: "sshd", PID: 500},
		},
	}
}

func TestBuildRollupReport_Counts(t *testing.T) {
	res := makeRollupResult()
	r := BuildRollupReport(res)
	if r.TotalPorts != 3 {
		t.Errorf("expected TotalPorts=3, got %d", r.TotalPorts)
	}
	if r.SnapshotCount != 5 {
		t.Errorf("expected SnapshotCount=5, got %d", r.SnapshotCount)
	}
	if r.Dropped != 2 {
		t.Errorf("expected Dropped=2, got %d", r.Dropped)
	}
}

func TestBuildRollupReport_TimestampSet(t *testing.T) {
	r := BuildRollupReport(makeRollupResult())
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestWriteRollupText_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	WriteRollupText(&buf, BuildRollupReport(makeRollupResult()))
	out := buf.String()
	for _, want := range []string{"Rollup Report", "PORT", "PROTOCOL", "PROCESS", "PID"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteRollupText_ContainsPortNumbers(t *testing.T) {
	var buf bytes.Buffer
	WriteRollupText(&buf, BuildRollupReport(makeRollupResult()))
	out := buf.String()
	for _, want := range []string{"80", "443", "22", "nginx", "sshd"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in text output", want)
		}
	}
}

func TestWriteRollupJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRollupJSON(&buf, BuildRollupReport(makeRollupResult())); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out RollupReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalPorts != 3 {
		t.Errorf("expected TotalPorts=3 in JSON, got %d", out.TotalPorts)
	}
}

func TestBuildRollupReport_EmptyInput(t *testing.T) {
	res := snapshot.RollupResult{Timestamp: time.Now()}
	r := BuildRollupReport(res)
	if r.TotalPorts != 0 {
		t.Errorf("expected 0 ports for empty input")
	}
	if len(r.Ports) != 0 {
		t.Errorf("expected empty ports slice")
	}
}
