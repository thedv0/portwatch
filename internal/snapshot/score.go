package snapshot

import (
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// RiskLevel represents a scored risk tier for an open port.
type RiskLevel int

const (
	RiskLow    RiskLevel = 0
	RiskMedium RiskLevel = 1
	RiskHigh   RiskLevel = 2
)

// ScoredPort pairs a scanned port with its computed risk level.
type ScoredPort struct {
	Port  scanner.Port
	Score RiskLevel
	Reason string
}

// ScoreOptions controls how risk scores are computed.
type ScoreOptions struct {
	// HighRiskPorts are port numbers always considered high risk.
	HighRiskPorts []int
	// MediumRiskPorts are port numbers considered medium risk.
	MediumRiskPorts []int
	// FlagPIDZero marks ports with PID 0 as high risk.
	FlagPIDZero bool
}

// DefaultScoreOptions returns a sensible default scoring configuration.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		HighRiskPorts:   []int{22, 23, 3389, 5900, 4444, 1337},
		MediumRiskPorts: []int{21, 25, 110, 143, 3306, 5432, 6379, 27017},
		FlagPIDZero:     true,
	}
}

// ScorePorts assigns a risk level to each port and returns sorted results.
// Results are sorted descending by score (high risk first).
func ScorePorts(ports []scanner.Port, opts ScoreOptions) []ScoredPort {
	high := toSet(opts.HighRiskPorts)
	med := toSet(opts.MediumRiskPorts)

	results := make([]ScoredPort, 0, len(ports))
	for _, p := range ports {
		sp := ScoredPort{Port: p, Score: RiskLow, Reason: "no known risk"}
		if opts.FlagPIDZero && p.PID == 0 {
			sp.Score = RiskHigh
			sp.Reason = "PID is zero (unknown process)"
		} else if high[p.Port] {
			sp.Score = RiskHigh
			sp.Reason = "high-risk port number"
		} else if med[p.Port] {
			sp.Score = RiskMedium
			sp.Reason = "medium-risk port number"
		}
		results = append(results, sp)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
