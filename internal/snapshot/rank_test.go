package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeRankedPort(port int, score int) ScoredPort {
	return ScoredPort{
		Port:  scanner.Port{Port: port, Protocol: "tcp"},
		Score: score,
	}
}

func TestRankPorts_Empty(t *testing.T) {
	result := RankPorts(nil, nil, DefaultRankOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestRankPorts_OrderedByScore(t *testing.T) {
	scored := []ScoredPort{
		makeRankedPort(80, 30),
		makeRankedPort(22, 90),
		makeRankedPort(443, 50),
	}
	freq := map[int]int{80: 5, 22: 1, 443: 3}
	opts := DefaultRankOptions()
	opts.TopN = 0

	result := RankPorts(scored, freq, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
	// port 22 has highest risk weight, should rank first
	if result[0].Port.Port != 22 {
		t.Errorf("expected port 22 first, got %d", result[0].Port.Port)
	}
}

func TestRankPorts_TopNLimitsResults(t *testing.T) {
	scored := []ScoredPort{
		makeRankedPort(22, 90),
		makeRankedPort(443, 50),
		makeRankedPort(80, 30),
		makeRankedPort(8080, 20),
	}
	opts := DefaultRankOptions()
	opts.TopN = 2

	result := RankPorts(scored, nil, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestRankPorts_RankFieldSet(t *testing.T) {
	scored := []ScoredPort{
		makeRankedPort(22, 80),
		makeRankedPort(80, 40),
	}
	opts := DefaultRankOptions()
	opts.TopN = 0

	result := RankPorts(scored, nil, opts)
	for i, r := range result {
		if r.Rank != i+1 {
			t.Errorf("rank[%d] = %d, want %d", i, r.Rank, i+1)
		}
	}
}

func TestRankPorts_FrequencyInfluencesScore(t *testing.T) {
	// Two ports with same risk; higher frequency should rank first.
	scored := []ScoredPort{
		makeRankedPort(80, 50),
		makeRankedPort(443, 50),
	}
	freq := map[int]int{80: 1, 443: 100}
	opts := RankOptions{TopN: 0, WeightRisk: 0.0, WeightFrequency: 1.0}

	result := RankPorts(scored, freq, opts)
	if result[0].Port.Port != 443 {
		t.Errorf("expected port 443 first due to frequency, got %d", result[0].Port.Port)
	}
}
