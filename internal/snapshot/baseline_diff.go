package snapshot

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// BaselineDiff describes deviations from a saved baseline.
type BaselineDiff struct {
	Baseline *Baseline
	Added    []scanner.Port // ports present now but not in baseline
	Removed  []scanner.Port // ports in baseline but not present now
}

// HasChanges returns true when there are any added or removed ports.
func (d *BaselineDiff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Summary returns a human-readable one-line description of the diff.
func (d *BaselineDiff) Summary() string {
	return fmt.Sprintf("baseline %q: +%d added, -%d removed",
		d.Baseline.Label, len(d.Added), len(d.Removed))
}

// CompareToBaseline diffs current ports against a stored baseline.
func CompareToBaseline(b *Baseline, current []scanner.Port) *BaselineDiff {
	baseIdx := indexPorts(b.Ports)
	currIdx := indexPorts(current)

	var added, removed []scanner.Port

	for key, p := range currIdx {
		if _, ok := baseIdx[key]; !ok {
			added = append(added, p)
		}
	}
	for key, p := range baseIdx {
		if _, ok := currIdx[key]; !ok {
			removed = append(removed, p)
		}
	}

	return &BaselineDiff{
		Baseline: b,
		Added:    added,
		Removed:  removed,
	}
}
