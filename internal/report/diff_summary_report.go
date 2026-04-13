package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// DiffSummaryReport wraps a DiffSummary for reporting.
type DiffSummaryReport struct {
	Timestamp    time.Time              `json:"timestamp"`
	Summary      snapshot.DiffSummary   `json:"summary"`
	HasChanges   bool                   `json:"has_changes"`
}

// BuildDiffSummaryReport constructs a report from a DiffResult.
func BuildDiffSummaryReport(d snapshot.DiffResult, clock func() time.Time) DiffSummaryReport {
	s := snapshot.SummarizeDiff(d, clock)
	return DiffSummaryReport{
		Timestamp:  s.Timestamp,
		Summary:    s,
		HasChanges: s.HasChanges,
	}
}

// WriteDiffSummaryText writes a human-readable diff summary to w.
func WriteDiffSummaryText(r DiffSummaryReport, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Diff Summary [%s]\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  Added:     %d\n", r.Summary.AddedCount)
	fmt.Fprintf(w, "  Removed:   %d\n", r.Summary.RemovedCount)
	fmt.Fprintf(w, "  Unchanged: %d\n", r.Summary.UnchangedCount)
	if len(r.Summary.AddedPorts) > 0 {
		fmt.Fprintf(w, "  New ports: %v\n", r.Summary.AddedPorts)
	}
	if len(r.Summary.RemovedPorts) > 0 {
		fmt.Fprintf(w, "  Gone ports: %v\n", r.Summary.RemovedPorts)
	}
	if !r.HasChanges {
		fmt.Fprintln(w, "  No changes detected.")
	}
	return nil
}

// WriteDiffSummaryJSON writes a JSON-encoded diff summary report to w.
func WriteDiffSummaryJSON(r DiffSummaryReport, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
