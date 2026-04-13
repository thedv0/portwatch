package snapshot

import (
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// RankOptions controls how ports are ranked.
type RankOptions struct {
	// TopN limits results to the top N ports. 0 means return all.
	TopN int
	// WeightRisk multiplies the risk score contribution.
	WeightRisk float64
	// WeightFrequency multiplies the frequency contribution.
	WeightFrequency float64
}

// DefaultRankOptions returns sensible defaults.
func DefaultRankOptions() RankOptions {
	return RankOptions{
		TopN:            10,
		WeightRisk:      0.7,
		WeightFrequency: 0.3,
	}
}

// RankedPort pairs a port with its computed rank score.
type RankedPort struct {
	Port  scanner.Port
	Score float64
	Rank  int
}

// RankPorts scores and orders ports by a weighted combination of risk and
// frequency derived from the provided scored ports and a frequency map.
// freqMap maps port number to observed frequency count.
func RankPorts(scored []ScoredPort, freqMap map[int]int, opts RankOptions) []RankedPort {
	if len(scored) == 0 {
		return nil
	}

	// Normalise frequency to [0,1].
	maxFreq := 1
	for _, f := range freqMap {
		if f > maxFreq {
			maxFreq = f
		}
	}

	ranked := make([]RankedPort, 0, len(scored))
	for _, sp := range scored {
		normRisk := float64(sp.Score) / 100.0
		normFreq := float64(freqMap[sp.Port.Port]) / float64(maxFreq)
		combined := opts.WeightRisk*normRisk + opts.WeightFrequency*normFreq
		ranked = append(ranked, RankedPort{
			Port:  sp.Port,
			Score: combined,
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	for i := range ranked {
		ranked[i].Rank = i + 1
	}

	if opts.TopN > 0 && len(ranked) > opts.TopN {
		ranked = ranked[:opts.TopN]
	}

	return ranked
}
