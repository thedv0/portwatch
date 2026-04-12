package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeTrendMetrics() map[string][]snapshot.TrendPoint {
	return map[string][]snapshot.TrendPoint{
		"open_ports": {
			{Timestamp: 1, Value: 10},
			{Timestamp: 2, Value: 12},
			{Timestamp: 3, Value: 14},
		},
		"listeners": {
			{Timestamp: 1, Value: 5},
			{Timestamp: 2, Value: 5},
			{Timestamp: 3, Value: 5},
		},
	}
}

func TestBuildTrendReport_HasEntries(t *testing.T) {
	r := BuildTrendReport(makeTrendMetrics(), snapshot.DefaultTrendOptions())
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestBuildTrendReport_DirectionCorrect(t *testing.T) {
	r := BuildTrendReport(makeTrendMetrics(), snapshot.DefaultTrendOptions())
	for _, e := range r.Entries {
		switch e.Metric {
		case "open_ports":
			if e.Direction != snapshot.TrendUp {
				t.Errorf("open_ports: expected up, got %s", e.Direction)
			}
		case "listeners":
			if e.Direction != snapshot.TrendStable {
				t.Errorf("listeners: expected stable, got %s", e.Direction)
			}
		}
	}
}

func TestWriteTrendText_ContainsHeaders(t *testing.T) {
	r := BuildTrendReport(makeTrendMetrics(), snapshot.DefaultTrendOptions())
	var buf bytes.Buffer
	if err := WriteTrendText(&buf, r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"METRIC", "DIRECTION", "SLOPE", "POINTS"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in text output", want)
		}
	}
}

func TestWriteTrendJSON_ValidJSON(t *testing.T) {
	r := BuildTrendReport(makeTrendMetrics(), snapshot.DefaultTrendOptions())
	var buf bytes.Buffer
	if err := WriteTrendJSON(&buf, r); err != nil {
		t.Fatal(err)
	}
	var out TrendReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries in JSON, got %d", len(out.Entries))
	}
}

func TestBuildTrendReport_EmptyMetrics(t *testing.T) {
	r := BuildTrendReport(map[string][]snapshot.TrendPoint{}, snapshot.DefaultTrendOptions())
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
	if r.GeneratedAt.IsZero() {
		t.Fatal("GeneratedAt should not be zero")
	}
}
