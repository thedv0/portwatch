package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeSplitResults() []snapshot.SplitResult {
	return []snapshot.SplitResult{
		{
			Index: 0,
			Ports: []snapshot.Port{
				{Port: 80, Protocol: "tcp", Process: "nginx", PID: 1},
				{Port: 443, Protocol: "tcp", Process: "nginx", PID: 1},
			},
		},
		{
			Index: 1,
			Ports: []snapshot.Port{
				{Port: 53, Protocol: "udp", Process: "dns", PID: 2},
			},
		},
	}
}

func TestBuildSplitReport_PartCount(t *testing.T) {
	res := makeSplitResults()
	r := BuildSplitReport(res, "protocol")
	if r.PartCount != 2 {
		t.Errorf("expected PartCount=2, got %d", r.PartCount)
	}
}

func TestBuildSplitReport_TotalPorts(t *testing.T) {
	res := makeSplitResults()
	r := BuildSplitReport(res, "protocol")
	if r.TotalPorts != 3 {
		t.Errorf("expected TotalPorts=3, got %d", r.TotalPorts)
	}
}

func TestBuildSplitReport_TimestampSet(t *testing.T) {
	res := makeSplitResults()
	r := BuildSplitReport(res, "port")
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestWriteSplitText_ContainsHeaders(t *testing.T) {
	res := makeSplitResults()
	r := BuildSplitReport(res, "protocol")
	var buf bytes.Buffer
	if err := WriteSplitText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Split Report", "Field", "Parts", "Total", "protocol"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteSplitText_ContainsPortNumbers(t *testing.T) {
	res := makeSplitResults()
	r := BuildSplitReport(res, "protocol")
	var buf bytes.Buffer
	_ = WriteSplitText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "80") || !strings.Contains(out, "443") {
		t.Error("expected port numbers in text output")
	}
}

func TestWriteSplitJSON_ValidJSON(t *testing.T) {
	res := makeSplitResults()
	r := BuildSplitReport(res, "port")
	var buf bytes.Buffer
	if err := WriteSplitJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["part_count"]; !ok {
		t.Error("expected part_count in JSON")
	}
}
