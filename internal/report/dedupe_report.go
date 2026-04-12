package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// DedupeReport summarises the result of a deduplication pass.
type DedupeReport struct {
	Timestamp  time.Time      `json:"timestamp"`
	Before     int            `json:"before"`
	After      int            `json:"after"`
	Removed    int            `json:"removed"`
	Resulting  []scanner.Port `json:"resulting"`
}

// BuildDedupeReport runs Dedupe and returns a populated DedupeReport.
func BuildDedupeReport(ports []scanner.Port, opts snapshot.DedupeOptions) DedupeReport {
	deduped := snapshot.Dedupe(ports, opts)
	return DedupeReport{
		Timestamp: time.Now().UTC(),
		Before:    len(ports),
		After:     len(deduped),
		Removed:   len(ports) - len(deduped),
		Resulting: deduped,
	}
}

// WriteDedupeText writes a human-readable dedupe report to w.
func WriteDedupeText(w io.Writer, r DedupeReport) error {
	_, err := fmt.Fprintf(w,
		"Dedupe Report (%s)\n  Before : %d\n  After  : %d\n  Removed: %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.Before, r.After, r.Removed,
	)
	if err != nil {
		return err
	}
	for _, p := range r.Resulting {
		if _, err := fmt.Fprintf(w, "  [%s] :%d pid=%d proc=%s\n",
			p.Protocol, p.Port, p.PID, p.Process); err != nil {
			return err
		}
	}
	return nil
}

// WriteDedupeJSON writes the dedupe report as JSON to w.
func WriteDedupeJSON(w io.Writer, r DedupeReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
