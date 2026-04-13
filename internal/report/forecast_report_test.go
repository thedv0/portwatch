package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeForecastResult() snapshot.ForecastResult {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return snapshot.ForecastResult{
		GeneratedAt: now,
		Horizon:     time.Hour,
		Slope:       0.05,
		Intercept:   10.0,
		Points: []snapshot.ForecastPoint{
			{At: now.Add(10 * time.Minute), PortCount: 10.3},
			{At: now.Add(20 * time.Minute), PortCount: 10.6},
		},
	}
}

func TestBuildForecastReport_FieldsPopulated(t *testing.T) {
	res := makeForecastResult()
	r := BuildForecastReport(res)
	if r.Slope != res.Slope {
		t.Errorf("slope mismatch: got %f want %f", r.Slope, res.Slope)
	}
	if r.Intercept != res.Intercept {
		t.Errorf("intercept mismatch")
	}
	if len(r.Points) != len(res.Points) {
		t.Errorf("point count mismatch: got %d want %d", len(r.Points), len(res.Points))
	}
	if r.Horizon != "1h0m0s" {
		t.Errorf("unexpected horizon string: %s", r.Horizon)
	}
}

func TestBuildForecastReport_TimestampSet(t *testing.T) {
	res := makeForecastResult()
	r := BuildForecastReport(res)
	if r.GeneratedAt.IsZero() {
		t.Error("generated_at should not be zero")
	}
}

func TestWriteForecastText_ContainsHeaders(t *testing.T) {
	r := BuildForecastReport(makeForecastResult())
	var buf bytes.Buffer
	if err := WriteForecastText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Forecast Report", "Slope", "Intercept", "Horizon", "Predicted Ports"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestWriteForecastText_ContainsPoints(t *testing.T) {
	r := BuildForecastReport(makeForecastResult())
	var buf bytes.Buffer
	_ = WriteForecastText(&buf, r)
	if !strings.Contains(buf.String(), "10.") {
		t.Error("expected port count values in text output")
	}
}

func TestWriteForecastJSON_ValidJSON(t *testing.T) {
	r := BuildForecastReport(makeForecastResult())
	var buf bytes.Buffer
	if err := WriteForecastJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["points"]; !ok {
		t.Error("JSON missing 'points' key")
	}
}
