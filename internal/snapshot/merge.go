package snapshot

import (
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// MergeStrategy controls how duplicate ports are handled when merging snapshots.
type MergeStrategy int

const (
	// MergeStrategyUnion keeps all unique ports from both snapshots.
	MergeStrategyUnion MergeStrategy = iota
	// MergeStrategyIntersect keeps only ports present in both snapshots.
	MergeStrategyIntersect
	// MergeStrategyPreferLeft keeps the left snapshot's entry on conflict.
	MergeStrategyPreferLeft
)

// MergeOptions configures the Merge operation.
type MergeOptions struct {
	Strategy MergeStrategy
}

// DefaultMergeOptions returns sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Strategy: MergeStrategyUnion,
	}
}

// Merge combines two port slices according to the given strategy.
// The result is sorted by port number then protocol.
func Merge(left, right []scanner.Port, opts MergeOptions) ([]scanner.Port, error) {
	switch opts.Strategy {
	case MergeStrategyUnion:
		return mergeUnion(left, right), nil
	case MergeStrategyIntersect:
		return mergeIntersect(left, right), nil
	case MergeStrategyPreferLeft:
		return mergePreferLeft(left, right), nil
	default:
		return nil, fmt.Errorf("unknown merge strategy: %d", opts.Strategy)
	}
}

func mergeKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Port)
}

func mergeUnion(left, right []scanner.Port) []scanner.Port {
	seen := make(map[string]scanner.Port)
	for _, p := range left {
		seen[mergeKey(p)] = p
	}
	for _, p := range right {
		if _, exists := seen[mergeKey(p)]; !exists {
			seen[mergeKey(p)] = p
		}
	}
	return sortedPorts(seen)
}

func mergeIntersect(left, right []scanner.Port) []scanner.Port {
	rightIdx := make(map[string]struct{})
	for _, p := range right {
		rightIdx[mergeKey(p)] = struct{}{}
	}
	seen := make(map[string]scanner.Port)
	for _, p := range left {
		if _, ok := rightIdx[mergeKey(p)]; ok {
			seen[mergeKey(p)] = p
		}
	}
	return sortedPorts(seen)
}

func mergePreferLeft(left, right []scanner.Port) []scanner.Port {
	seen := make(map[string]scanner.Port)
	for _, p := range right {
		seen[mergeKey(p)] = p
	}
	for _, p := range left {
		seen[mergeKey(p)] = p
	}
	return sortedPorts(seen)
}

func sortedPorts(m map[string]scanner.Port) []scanner.Port {
	out := make([]scanner.Port, 0, len(m))
	for _, p := range m {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Port != out[j].Port {
			return out[i].Port < out[j].Port
		}
		return out[i].Protocol < out[j].Protocol
	})
	return out
}
