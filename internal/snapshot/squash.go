package snapshot

import (
	"errors"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// SquashOptions controls how multiple snapshots are squashed into one.
type SquashOptions struct {
	// Strategy determines how port conflicts are resolved.
	// "first" keeps the earliest occurrence, "last" keeps the latest.
	Strategy string
	// Label is applied to the resulting squashed snapshot.
	Label string
}

// DefaultSquashOptions returns sensible defaults.
func DefaultSquashOptions() SquashOptions {
	return SquashOptions{
		Strategy: "last",
		Label:    "squashed",
	}
}

// SquashResult holds the output of a Squash operation.
type SquashResult struct {
	Ports     []scanner.Port
	Timestamp time.Time
	Label     string
	InputSnaps int
}

// Squash collapses multiple snapshots into a single deduplicated snapshot.
// Port identity is based on (protocol, port). Conflicts are resolved per opts.Strategy.
func Squash(snaps []Snapshot, opts SquashOptions) (SquashResult, error) {
	if opts.Strategy == "" {
		opts.Strategy = "last"
	}
	if opts.Strategy != "first" && opts.Strategy != "last" {
		return SquashResult{}, errors.New("squash: strategy must be \"first\" or \"last\"")
	}
	if len(snaps) == 0 {
		return SquashResult{
			Timestamp:  time.Now(),
			Label:      opts.Label,
			InputSnaps: 0,
		}, nil
	}

	type key struct {
		Proto string
		Port  int
	}

	seen := make(map[key]scanner.Port)

	for _, snap := range snaps {
		for _, p := range snap.Ports {
			k := key{Proto: p.Protocol, Port: p.Port}
			if _, exists := seen[k]; !exists || opts.Strategy == "last" {
				seen[k] = p
			}
		}
	}

	result := make([]scanner.Port, 0, len(seen))
	for _, p := range seen {
		result = append(result, p)
	}

	ts := snaps[len(snaps)-1].Timestamp
	if opts.Strategy == "first" {
		ts = snaps[0].Timestamp
	}

	return SquashResult{
		Ports:      result,
		Timestamp:  ts,
		Label:      opts.Label,
		InputSnaps: len(snaps),
	}, nil
}
