package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/netwatch/portwatch/internal/scanner"
	"github.com/netwatch/portwatch/internal/snapshot"
)

func makePipelineResult(inputLen, outputLen, issues int) (int, snapshot.PipelineResult) {
	ports := make([]scanner.Port, outputLen)
	for i := range ports {
		ports[i] = scanner.Port{Port: 8000 + i, Protocol: "tcp", Process: "svc", PID: i + 1}
	}
	vr := snapshot.ValidationResult{Issues: make([]snapshot.ValidationIssue, issues)}
	return inputLen, snapshot.PipelineResult{
		Ports:      ports,
		Validation: &vr,
		Classified: make([]snapshot.ClassifiedPort, outputLen),
	}
}

func TestBuildPipelineReport_Counts(t *testing.T) {
	inputLen, result := makePipelineResult(5, 3, 2)
	r := BuildPipelineReport(inputLen, result)

	if r.TotalInput != 5 {
		t.Errorf("expected TotalInput=5, got %d", r.TotalInput)
	}
	if r.TotalOutput != 3 {
		t.Errorf("expected TotalOutput=3, got %d", r.TotalOutput)
	}
	if r.Duplicates != 2 {
		t.Errorf("expected Duplicates=2, got %d", r.Duplicates)
	}
	if r.Issues != 2 {
		t.Errorf("expected Issues=2, got %d", r.Issues)
	}
}

func TestBuildPipelineReport_NilValidation(t *testing.T) {
	result := snapshot.PipelineResult{Ports: []scanner.Port{}}
	r := BuildPipelineReport(0, result)
	if r.Issues != 0 {
		t.Errorf("expected 0 issues for nil validation, got %d", r.Issues)
	}
}

func TestWritePipelineText_ContainsHeaders(t *testing.T) {
	_, result := makePipelineResult(4, 4, 0)
	r := BuildPipelineReport(4, result)
	var buf bytes.Buffer
	if err := WritePipelineText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Pipeline Report", "Input", "Output", "Duplicates", "Issues", "Classified"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWritePipelineJSON_ValidJSON(t *testing.T) {
	_, result := makePipelineResult(3, 2, 1)
	r := BuildPipelineReport(3, result)
	var buf bytes.Buffer
	if err := WritePipelineJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out PipelineReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalInput != 3 {
		t.Errorf("expected TotalInput=3, got %d", out.TotalInput)
	}
}

func TestBuildPipelineReport_TimestampSet(t *testing.T) {
	_, result := makePipelineResult(1, 1, 0)
	r := BuildPipelineReport(1, result)
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
