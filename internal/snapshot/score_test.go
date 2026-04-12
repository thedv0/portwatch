package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeScorePort(port, pid int) scanner.Port {
	return scanner.Port{Port: port, PID: pid, Protocol: "tcp"}
}

func TestScorePorts_HighRiskPort(t *testing.T) {
	opts := DefaultScoreOptions()
	ports := []scanner.Port{makeScorePort(22, 1000)}
	results := ScorePorts(ports, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Score != RiskHigh {
		t.Errorf("expected RiskHigh for port 22, got %v", results[0].Score)
	}
}

func TestScorePorts_MediumRiskPort(t *testing.T) {
	opts := DefaultScoreOptions()
	ports := []scanner.Port{makeScorePort(3306, 500)}
	results := ScorePorts(ports, opts)
	if results[0].Score != RiskMedium {
		t.Errorf("expected RiskMedium for port 3306, got %v", results[0].Score)
	}
}

func TestScorePorts_LowRiskPort(t *testing.T) {
	opts := DefaultScoreOptions()
	ports := []scanner.Port{makeScorePort(8080, 1234)}
	results := ScorePorts(ports, opts)
	if results[0].Score != RiskLow {
		t.Errorf("expected RiskLow for port 8080, got %v", results[0].Score)
	}
}

func TestScorePorts_PIDZeroIsHighRisk(t *testing.T) {
	opts := DefaultScoreOptions()
	ports := []scanner.Port{makeScorePort(9999, 0)}
	results := ScorePorts(ports, opts)
	if results[0].Score != RiskHigh {
		t.Errorf("expected RiskHigh for PID 0, got %v", results[0].Score)
	}
}

func TestScorePorts_SortedDescending(t *testing.T) {
	opts := DefaultScoreOptions()
	ports := []scanner.Port{
		makeScorePort(8080, 100),  // low
		makeScorePort(3306, 200),  // medium
		makeScorePort(22, 300),    // high
	}
	results := ScorePorts(ports, opts)
	if results[0].Score != RiskHigh || results[1].Score != RiskMedium || results[2].Score != RiskLow {
		t.Errorf("results not sorted descending: %v %v %v",
			results[0].Score, results[1].Score, results[2].Score)
	}
}

func TestScorePorts_FlagPIDZeroDisabled(t *testing.T) {
	opts := DefaultScoreOptions()
	opts.FlagPIDZero = false
	ports := []scanner.Port{makeScorePort(9999, 0)}
	results := ScorePorts(ports, opts)
	if results[0].Score != RiskLow {
		t.Errorf("expected RiskLow when FlagPIDZero disabled, got %v", results[0].Score)
	}
}

func TestScorePorts_EmptyInput(t *testing.T) {
	opts := DefaultScoreOptions()
	results := ScorePorts(nil, opts)
	if len(results) != 0 {
		t.Errorf("expected empty results for nil input")
	}
}

func TestDefaultScoreOptions_HasHighRiskPorts(t *testing.T) {
	opts := DefaultScoreOptions()
	if len(opts.HighRiskPorts) == 0 {
		t.Error("expected non-empty HighRiskPorts")
	}
	if !opts.FlagPIDZero {
		t.Error("expected FlagPIDZero to be true by default")
	}
}
