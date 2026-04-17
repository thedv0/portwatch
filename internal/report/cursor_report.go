package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// CursorReport summarises a cursor page result.
type CursorReport struct {
	Timestamp  time.Time          `json:"timestamp"`
	PageSize   int                `json:"page_size"`
	Returned   int                `json:"returned"`
	HasMore    bool               `json:"has_more"`
	NextAfter  time.Time          `json:"next_after"`
	Snapshots  []snapshot.Snapshot `json:"snapshots"`
}

// BuildCursorReport constructs a CursorReport from a CursorResult.
func BuildCursorReport(res snapshot.CursorResult) CursorReport {
	return CursorReport{
		Timestamp: time.Now().UTC(),
		PageSize:  len(res.Page),
		Returned:  len(res.Page),
		HasMore:   res.HasMore,
		NextAfter: res.NextAfter,
		Snapshots: res.Page,
	}
}

// WriteCursorText writes a human-readable cursor report to w.
func WriteCursorText(w io.Writer, r CursorReport) error {
	fmt.Fprintf(w, "Cursor Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Returned : %d\n", r.Returned)
	fmt.Fprintf(w, "Has More : %v\n", r.HasMore)
	if !r.NextAfter.IsZero() {
		fmt.Fprintf(w, "Next After: %s\n", r.NextAfter.Format(time.RFC3339))
	}
	fmt.Fprintln(w, "---")
	for _, s := range r.Snapshots {
		fmt.Fprintf(w, "  [%s] ports=%d\n", s.Timestamp.Format(time.RFC3339), len(s.Ports))
	}
	return nil
}

// WriteCursorJSON writes the report as JSON to w.
func WriteCursorJSON(w io.Writer, r CursorReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
