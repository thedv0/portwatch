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

func makeScored(port, pid int, level snapshot.RiskLevel, reason string) snapshot.ScoredPort {
	return snapshot.ScoredPort{
		Port:   scanner.Port{Port: port, PID: pid, Protocol: "tcp"},
		Score:  level,
		Reason: reason,
	}
}

func TestBuildScoreReport_Partitions(t *testing.T) {
	scored := []snapshot.ScoredPort{
		makeScored(22, 100, snapshot.RiskHigh, "high-risk port"),
		makeScored(3306, 200, snapshot.RiskMedium, "medium-risk port"),
		makeScored(8080, 300, snapshot.RiskLow, "no known risk"),
	}
	r := BuildScoreReport(scored, time.Now())
	if r.TotalPorts != 3 {
		t.Errorf("expected TotalPorts=3, got %d", r.TotalPorts)
	}
	if len(r.HighRisk) != 1 || len(r.MediumRisk) != 1 || len(r.LowRisk) != 1 {
		t.Errorf("unexpected tier counts: high=%d med=%d low=%d",
			len(r.HighRisk), len(r.MediumRisk), len(r.LowRisk))
	}
}

func TestWriteScoreText_ContainsRiskLabels(t *testing.T) {
	scored := []snapshot.ScoredPort{
		makeScored(22, 1, snapshot.RiskHigh, "high-risk port"),
		makeScored(8080, 2, snapshot.RiskLow, "no known risk"),
	}
	r := BuildScoreReport(scored, time.Now())
	var buf bytes.Buffer
	if err := WriteScoreText(&buf, r); err != nil {
		t.Fatalf("WriteScoreText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[HIGH]") {
		t.Error("expected [HIGH] in output")
	}
	if !strings.Contains(out, "[LOW]") {
		t.Error("expected [LOW] in output")
	}
	if !strings.Contains(out, "high-risk port") {
		t.Error("expected reason in output")
	}
}

func TestWriteScoreJSON_ValidJSON(t *testing.T) {
	scored := []snapshot.ScoredPort{
		makeScored(22, 1, snapshot.RiskHigh, "high-risk port"),
	}
	r := BuildScoreReport(scored, time.Now())
	var buf bytes.Buffer
	if err := WriteScoreJSON(&buf, r); err != nil {
		t.Fatalf("WriteScoreJSON error: %v", err)
	}
	var out ScoreReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalPorts != 1 {
		t.Errorf("expected TotalPorts=1, got %d", out.TotalPorts)
	}
}

func TestBuildScoreReport_EmptyInput(t *testing.T) {
	r := BuildScoreReport(nil, time.Now())
	if r.TotalPorts != 0 {
		t.Errorf("expected TotalPorts=0, got %d", r.TotalPorts)
	}
	if len(r.HighRisk) != 0 || len(r.MediumRisk) != 0 || len(r.LowRisk) != 0 {
		t.Error("expected all tiers empty for empty input")
	}
}
