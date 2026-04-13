package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEvictPorts(ports ...int) []snapshot.PortState {
	out := make([]snapshot.PortState, len(ports))
	for i, p := range ports {
		out[i] = snapshot.PortState{Port: p, Protocol: "tcp"}
	}
	return out
}

func TestBuildEvictReport_Counts(t *testing.T) {
	before := makeEvictPorts(80, 443, 8080)
	after := makeEvictPorts(443, 8080)
	r := BuildEvictReport("age", before, after)

	if r.Before != 3 {
		t.Errorf("Before: want 3, got %d", r.Before)
	}
	if r.After != 2 {
		t.Errorf("After: want 2, got %d", r.After)
	}
	if r.Evicted != 1 {
		t.Errorf("Evicted: want 1, got %d", r.Evicted)
	}
}

func TestBuildEvictReport_Policy(t *testing.T) {
	r := BuildEvictReport("count", makeEvictPorts(80), makeEvictPorts(80))
	if r.Policy != "count" {
		t.Errorf("Policy: want count, got %s", r.Policy)
	}
}

func TestBuildEvictReport_TimestampSet(t *testing.T) {
	before := time.Now()
	r := BuildEvictReport("idle", nil, nil)
	if r.Timestamp.Before(before) {
		t.Error("Timestamp should be set to approximately now")
	}
}

func TestBuildEvictReport_EmptyInput(t *testing.T) {
	r := BuildEvictReport("age", nil, nil)
	if r.Evicted != 0 {
		t.Errorf("expected 0 evicted, got %d", r.Evicted)
	}
}

func TestWriteEvictText_ContainsHeaders(t *testing.T) {
	r := BuildEvictReport("age", makeEvictPorts(80, 443), makeEvictPorts(443))
	var buf bytes.Buffer
	if err := WriteEvictText(&buf, r); err != nil {
		t.Fatalf("WriteEvictText error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Evict Report", "Policy", "Before", "After", "Evicted"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestWriteEvictText_ContainsRetained(t *testing.T) {
	r := BuildEvictReport("age", makeEvictPorts(80, 443), makeEvictPorts(443))
	var buf bytes.Buffer
	_ = WriteEvictText(&buf, r)
	if !strings.Contains(buf.String(), "443") {
		t.Error("output should contain retained port 443")
	}
}

func TestWriteEvictJSON_ValidJSON(t *testing.T) {
	r := BuildEvictReport("count", makeEvictPorts(80), makeEvictPorts(80))
	var buf bytes.Buffer
	if err := WriteEvictJSON(&buf, r); err != nil {
		t.Fatalf("WriteEvictJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["evicted"]; !ok {
		t.Error("JSON missing 'evicted' field")
	}
}
