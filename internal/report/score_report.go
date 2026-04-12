package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ScoreReport holds a scored snapshot summary for reporting.
type ScoreReport struct {
	Timestamp  time.Time              `json:"timestamp"`
	TotalPorts int                    `json:"total_ports"`
	HighRisk   []snapshot.ScoredPort  `json:"high_risk"`
	MediumRisk []snapshot.ScoredPort  `json:"medium_risk"`
	LowRisk    []snapshot.ScoredPort  `json:"low_risk"`
}

// BuildScoreReport partitions scored ports into risk tiers.
func BuildScoreReport(scored []snapshot.ScoredPort, ts time.Time) ScoreReport {
	r := ScoreReport{Timestamp: ts, TotalPorts: len(scored)}
	for _, sp := range scored {
		switch sp.Score {
		case snapshot.RiskHigh:
			r.HighRisk = append(r.HighRisk, sp)
		case snapshot.RiskMedium:
			r.MediumRisk = append(r.MediumRisk, sp)
		default:
			r.LowRisk = append(r.LowRisk, sp)
		}
	}
	return r
}

// WriteScoreText writes a human-readable risk report to w.
func WriteScoreText(w io.Writer, r ScoreReport) error {
	_, err := fmt.Fprintf(w, "Port Risk Report — %s\n", r.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Total: %d  High: %d  Medium: %d  Low: %d\n",
		r.TotalPorts, len(r.HighRisk), len(r.MediumRisk), len(r.LowRisk))
	for _, tier := range []struct {
		label string
		ports []snapshot.ScoredPort
	}{
		{"HIGH", r.HighRisk},
		{"MEDIUM", r.MediumRisk},
		{"LOW", r.LowRisk},
	} {
		for _, sp := range tier.ports {
			fmt.Fprintf(w, "  [%s] %s:%d pid=%d — %s\n",
				tier.label, sp.Port.Protocol, sp.Port.Port, sp.Port.PID, sp.Reason)
		}
	}
	return nil
}

// WriteScoreJSON writes the report as JSON to w.
func WriteScoreJSON(w io.Writer, r ScoreReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
