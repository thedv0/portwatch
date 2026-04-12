package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeEnrichedPorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Protocol: "tcp", Process: "ssh", PID: 1},
		{Port: 80, Protocol: "tcp", Process: "http", PID: 200},
		{Port: 9999, Protocol: "udp", Process: "", PID: 0},
	}
}

func TestBuildEnrichReport_Total(t *testing.T) {
	ports := makeEnrichedPorts()
	r := BuildEnrichReport(ports)
	if r.Total != 3 {
		t.Errorf("expected total 3, got %d", r.Total)
	}
}

func TestBuildEnrichReport_TimestampSet(t *testing.T) {
	r := BuildEnrichReport(makeEnrichedPorts())
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestWriteEnrichText_ContainsPortNumbers(t *testing.T) {
	var buf bytes.Buffer
	r := BuildEnrichReport(makeEnrichedPorts())
	if err := WriteEnrichText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	for _, want := range []string{"22", "80", "9999"} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestWriteEnrichText_UnknownProcessLabel(t *testing.T) {
	var buf bytes.Buffer
	r := BuildEnrichReport([]scanner.Port{{Port: 9999, Protocol: "udp", Process: "", PID: 0}})
	_ = WriteEnrichText(&buf, r)
	if !strings.Contains(buf.String(), "(unknown)") {
		t.Error("expected '(unknown)' for empty process")
	}
}

func TestWriteEnrichJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := BuildEnrichReport(makeEnrichedPorts())
	if err := WriteEnrichJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out EnrichReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Total != 3 {
		t.Errorf("expected total 3 in JSON, got %d", out.Total)
	}
}

func TestWriteEnrichJSON_PortsPreserved(t *testing.T) {
	var buf bytes.Buffer
	r := BuildEnrichReport(makeEnrichedPorts())
	_ = WriteEnrichJSON(&buf, r)
	var out EnrichReport
	_ = json.Unmarshal(buf.Bytes(), &out)
	if len(out.Ports) != 3 {
		t.Errorf("expected 3 ports in JSON output, got %d", len(out.Ports))
	}
}
