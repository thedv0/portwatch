package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeGroupReport() GroupReport {
	return GroupReport{
		GroupBy: "protocol",
		Groups: []snapshot.Group{
			{
				Key: "tcp",
				Ports: []scanner.Port{
					{Protocol: "tcp", Port: 80, PID: 100, Process: "nginx"},
					{Protocol: "tcp", Port: 443, PID: 100, Process: "nginx"},
				},
			},
			{
				Key: "udp",
				Ports: []scanner.Port{
					{Protocol: "udp", Port: 53, PID: 200, Process: "dns"},
				},
			},
		},
	}
}

func TestWriteGroupText_ContainsGroupKeys(t *testing.T) {
	var buf bytes.Buffer
	r := makeGroupReport()
	if err := WriteGroupText(&buf, r); err != nil {
		t.Fatalf("WriteGroupText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[tcp]") {
		t.Errorf("expected [tcp] in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[udp]") {
		t.Errorf("expected [udp] in output, got:\n%s", out)
	}
}

func TestWriteGroupText_ContainsPortNumbers(t *testing.T) {
	var buf bytes.Buffer
	r := makeGroupReport()
	_ = WriteGroupText(&buf, r)
	out := buf.String()
	for _, want := range []string{"80", "443", "53"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected port %s in output", want)
		}
	}
}

func TestWriteGroupText_ContainsProcessNames(t *testing.T) {
	var buf bytes.Buffer
	r := makeGroupReport()
	if err := WriteGroupText(&buf, r); err != nil {
		t.Fatalf("WriteGroupText error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"nginx", "dns"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected process name %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteGroupJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := makeGroupReport()
	if err := WriteGroupJSON(&buf, r); err != nil {
		t.Fatalf("WriteGroupJSON error: %v", err)
	}
	var out GroupReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.GroupBy != "protocol" {
		t.Errorf("expected group_by=protocol, got %s", out.GroupBy)
	}
	if len(out.Groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(out.Groups))
	}
}

func TestWriteGroupJSON_GroupPortCount(t *testing.T) {
	var buf bytes.Buffer
	r := makeGroupReport()
	_ = WriteGroupJSON(&buf, r)
	var out GroupReport
	_ = json.Unmarshal(buf.Bytes(), &out)
	if len(out.Groups[0].Ports) != 2 {
		t.Errorf("expected 2 ports in tcp group, got %d", len(out.Groups[0].Ports))
	}
}
