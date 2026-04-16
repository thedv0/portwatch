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

func makeUnionPorts() []scanner.Port {
	return []scanner.Port{
		{Protocol: "tcp", Port: 80, PID: 1, Process: "nginx"},
		{Protocol: "tcp", Port: 443, PID: 2, Process: "nginx"},
	}
}

func TestBuildUnionReport_Counts(t *testing.T) {
	ports := makeUnionPorts()
	snaps := []snapshot.Snapshot{{}, {}}
	opts := snapshot.DefaultUnionOptions()
	r := BuildUnionReport(snaps, ports, opts)
	if r.InputSnaps != 2 {
		t.Errorf("expected InputSnaps=2, got %d", r.InputSnaps)
	}
	if r.TotalPorts != 2 {
		t.Errorf("expected TotalPorts=2, got %d", r.TotalPorts)
	}
}

func TestBuildUnionReport_TimestampSet(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	r := BuildUnionReport(nil, nil, snapshot.DefaultUnionOptions())
	if r.Timestamp.Before(before) {
		t.Error("timestamp should be recent")
	}
}

func TestBuildUnionReport_OptionsPreserved(t *testing.T) {
	opts := snapshot.DefaultUnionOptions()
	opts.Dedup = false
	r := BuildUnionReport(nil, nil, opts)
	if r.DedupEnabled {
		t.Error("expected DedupEnabled=false")
	}
}

func TestWriteUnionText_ContainsHeaders(t *testing.T) {
	r := BuildUnionReport(
		[]snapshot.Snapshot{{}},
		makeUnionPorts(),
		snapshot.DefaultUnionOptions(),
	)
	var buf bytes.Buffer
	if err := WriteUnionText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Union Report", "Input Snaps", "Total Ports", "Dedup", "nginx"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteUnionJSON_ValidJSON(t *testing.T) {
	r := BuildUnionReport(
		[]snapshot.Snapshot{{}},
		makeUnionPorts(),
		snapshot.DefaultUnionOptions(),
	)
	var buf bytes.Buffer
	if err := WriteUnionJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out UnionReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalPorts != 2 {
		t.Errorf("expected TotalPorts=2, got %d", out.TotalPorts)
	}
}
