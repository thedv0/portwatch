package snapshot

import "time"

// DiffSummary holds a human-readable summary of changes between two snapshots.
type DiffSummary struct {
	Timestamp   time.Time `json:"timestamp"`
	AddedCount  int       `json:"added_count"`
	RemovedCount int      `json:"removed_count"`
	UnchangedCount int   `json:"unchanged_count"`
	AddedPorts  []string  `json:"added_ports,omitempty"`
	RemovedPorts []string `json:"removed_ports,omitempty"`
	HasChanges  bool      `json:"has_changes"`
}

// SummarizeDiff builds a DiffSummary from a DiffResult.
func SummarizeDiff(d DiffResult, clock func() time.Time) DiffSummary {
	if clock == nil {
		clock = time.Now
	}

	added := make([]string, 0, len(d.Added))
	for _, p := range d.Added {
		added = append(added, portKey(p.Protocol, itoa(p.Port)))
	}

	removed := make([]string, 0, len(d.Removed))
	for _, p := range d.Removed {
		removed = append(removed, portKey(p.Protocol, itoa(p.Port)))
	}

	return DiffSummary{
		Timestamp:      clock(),
		AddedCount:     len(d.Added),
		RemovedCount:   len(d.Removed),
		UnchangedCount: len(d.Unchanged),
		AddedPorts:     added,
		RemovedPorts:   removed,
		HasChanges:     len(d.Added) > 0 || len(d.Removed) > 0,
	}
}
