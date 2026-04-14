package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeProjectResult(steps int) snapshot.ProjectResult {
	now := time.Now()
	points := make([]snapshot.ProjectedPoint, steps)
	for i := range points {
		points[i] = snapshot.ProjectedPoint{
			At:    now.Add(time.Duration(i+1) * time.Minute),
			Count: float64(10 + i*2),
		}
	}
	return snapshot.ProjectResult{
		Points:      points,
		BaseCount:   10,
		Slope:       2.0,
		GeneratedAt: now,
	}
}

func TestBuildProjectReport_FieldsPopulated(t *testing.T) {
	res := makeProjectResult(3)
	r := BuildProjectReport(res)
	if r.Steps != 3 {
		t.Errorf("expected Steps=3, got %d", r.Steps)
	}
	if r.BaseCount != 10 {
		t.Errorf("expected BaseCount=10, got %f", r.BaseCount)
	}
	if r.Slope != 2.0 {
		t.Errorf("expected Slope=2.0, got %f", r.Slope)
	}
}

func TestBuildProjectReport_TimestampSet(t *testing.T) {
	res := makeProjectResult(2)
	r := BuildProjectReport(res)
	if r.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestWriteProjectText_ContainsHeaders(t *testing.T) {
	r := BuildProjectReport(makeProjectResult(2))
	var buf bytes.Buffer
	if err := WriteProjectText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Projection Report", "Base open ports", "Trend slope", "Steps projected", "Projected Count"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteProjectJSON_ValidJSON(t *testing.T) {
	r := BuildProjectReport(makeProjectResult(3))
	var buf bytes.Buffer
	if err := WriteProjectJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded ProjectReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Steps != 3 {
		t.Errorf("expected Steps=3 in decoded JSON, got %d", decoded.Steps)
	}
}

func TestWriteProjectText_ContainsPoints(t *testing.T) {
	r := BuildProjectReport(makeProjectResult(2))
	var buf bytes.Buffer
	_ = WriteProjectText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "10.00") && !strings.Contains(out, "12.00") {
		t.Error("expected projected counts in text output")
	}
}
