package snapshot

import (
	"errors"
	"time"
)

// ChainEntry holds a snapshot and its diff relative to the previous entry.
type ChainEntry struct {
	Snap Snapshot
	Diff DiffResult
	Index int
}

// DiffChainOptions controls how a diff chain is built.
type DiffChainOptions struct {
	// MaxEntries caps the number of chain entries (0 = unlimited).
	MaxEntries int
	// Since filters snapshots older than this time (zero = no filter).
	Since time.Time
}

// DefaultDiffChainOptions returns sensible defaults.
func DefaultDiffChainOptions() DiffChainOptions {
	return DiffChainOptions{MaxEntries: 0}
}

// BuildDiffChain computes a sequential diff chain across a slice of snapshots.
// Each entry contains the diff from the previous snapshot to the current one.
// The first entry always has an empty diff.
func BuildDiffChain(snaps []Snapshot, opts DiffChainOptions) ([]ChainEntry, error) {
	if len(snaps) == 0 {
		return nil, errors.New("diff_chain: no snapshots provided")
	}

	filtered := snaps
	if !opts.Since.IsZero() {
		filtered = make([]Snapshot, 0, len(snaps))
		for _, s := range snaps {
			if !s.Timestamp.Before(opts.Since) {
				filtered = append(filtered, s)
			}
		}
	}

	if opts.MaxEntries > 0 && len(filtered) > opts.MaxEntries {
		filtered = filtered[len(filtered)-opts.MaxEntries:]
	}

	chain := make([]ChainEntry, 0, len(filtered))
	for i, snap := range filtered {
		entry := ChainEntry{Snap: snap, Index: i}
		if i == 0 {
			entry.Diff = DiffResult{}
		} else {
			entry.Diff = Diff(filtered[i-1].Ports, snap.Ports)
		}
		chain = append(chain, entry)
	}
	return chain, nil
}
