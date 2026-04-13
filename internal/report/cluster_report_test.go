package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/netwatch/portwatch/internal/scanner"
)

func makeClusterMap() map[string][]scanner.Port {
	return map[string][]scanner.Port{
		"nginx": {
			{Port: 80, Protocol: "tcp", Process: "nginx", PID: 10},
			{Port: 443, Protocol: "tcp", Process: "nginx", PID: 10},
		},
		"postgres": {
			{Port: 5432, Protocol: "tcp", Process: "postgres", PID: 20},
		},
	}
}

func TestBuildClusterReport_ClusterCount(t *testing.T) {
	r := BuildClusterReport(makeClusterMap(), "process")
	if r.ClusterCount != 2 {
		t.Errorf("expected 2 clusters, got %d", r.ClusterCount)
	}
}

func TestBuildClusterReport_TotalPorts(t *testing.T) {
	r := BuildClusterReport(makeClusterMap(), "process")
	if r.TotalPorts != 3 {
		t.Errorf("expected 3 total ports, got %d", r.TotalPorts)
	}
}

func TestBuildClusterReport_TimestampSet(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	r := BuildClusterReport(makeClusterMap(), "process")
	if r.Timestamp.Before(before) {
		t.Error("timestamp should be recent")
	}
}

func TestWriteClusterText_ContainsGroupKey(t *testing.T) {
	r := BuildClusterReport(makeClusterMap(), "process")
	var buf bytes.Buffer
	if err := WriteClusterText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "nginx") {
		t.Error("expected 'nginx' in output")
	}
	if !strings.Contains(out, "postgres") {
		t.Error("expected 'postgres' in output")
	}
}

func TestWriteClusterText_ContainsPortNumbers(t *testing.T) {
	r := BuildClusterReport(makeClusterMap(), "process")
	var buf bytes.Buffer
	_ = WriteClusterText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "80") {
		t.Error("expected port 80 in output")
	}
	if !strings.Contains(out, "5432") {
		t.Error("expected port 5432 in output")
	}
}

func TestWriteClusterJSON_ValidJSON(t *testing.T) {
	r := BuildClusterReport(makeClusterMap(), "process")
	var buf bytes.Buffer
	if err := WriteClusterJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["clusters"]; !ok {
		t.Error("expected 'clusters' key in JSON output")
	}
}

func TestBuildClusterReport_EmptyInput(t *testing.T) {
	r := BuildClusterReport(nil, "process")
	if r.ClusterCount != 0 || r.TotalPorts != 0 {
		t.Error("expected zero counts for empty input")
	}
}
