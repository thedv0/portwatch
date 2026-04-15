package snapshot

import "time"

// DedupeChainEntry holds a snapshot and the number of duplicates removed.
type DedupeChainEntry struct {
	Snapshot  Snapshot
	Removed   int
	Timestamp time.Time
}

// DedupeChainOptions configures BuildDedupeChain.
type DedupeChainOptions struct {
	Dedupe DedupeOptions
}

// DefaultDedupeChainOptions returns sensible defaults.
func DefaultDedupeChainOptions() DedupeChainOptions {
	return DedupeChainOptions{
		Dedupe: DefaultDedupeOptions(),
	}
}

// BuildDedupeChain applies Dedupe to each snapshot in the slice and returns
// a chain of entries recording how many ports were removed per snapshot.
func BuildDedupeChain(snaps []Snapshot, opts DedupeChainOptions) []DedupeChainEntry {
	if len(snaps) == 0 {
		return nil
	}
	entries := make([]DedupeChainEntry, 0, len(snaps))
	for _, s := range snaps {
		before := len(s.Ports)
		deduped := Dedupe(s.Ports, opts.Dedupe)
		after := len(deduped)
		result := s
		result.Ports = deduped
		entries = append(entries, DedupeChainEntry{
			Snapshot:  result,
			Removed:   before - after,
			Timestamp: s.Timestamp,
		})
	}
	return entries
}
