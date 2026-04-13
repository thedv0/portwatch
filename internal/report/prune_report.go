package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// PruneReport summarises the result of a prune operation.
type PruneReport struct {
	Timestamp   time.Time            `json:"timestamp"`
	Before      int                  `json:"before"`
	After       int                  `json:"after"`
	Removed     int                  `json:"removed"`
	Retained    []snapshot.PortState `json:"retained"`
}

// BuildPruneReport constructs a PruneReport given the original and pruned slices.
func BuildPruneReport(before, after []snapshot.PortState) PruneReport {
	return PruneReport{
		Timestamp: time.Now().UTC(),
		Before:    len(before),
		After:     len(after),
		Removed:   len(before) - len(after),
		Retained:  after,
	}
}

// WritePruneText writes a human-readable prune summary to w.
func WritePruneText(w io.Writer, r PruneReport) error {
	_, err := fmt.Fprintf(w,
		"Prune Report [%s]\n  Before : %d\n  After  : %d\n  Removed: %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.Before,
		r.After,
		r.Removed,
	)
	if err != nil {
		return err
	}
	if len(r.Retained) == 0 {
		_, err = fmt.Fprintln(w, "  (no ports retained)")
		return err
	}
	for _, p := range r.Retained {
		_, err = fmt.Fprintf(w, "  %-6s %d\n", p.Protocol, p.Port)
		if err != nil {
			return err
		}
	}
	return nil
}

// WritePruneJSON serialises r as JSON into w.
func WritePruneJSON(w io.Writer, r PruneReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
