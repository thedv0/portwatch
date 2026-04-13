package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/wricardo/portwatch/internal/scanner"
	"github.com/wricardo/portwatch/internal/snapshot"
)

// SliceReport holds metadata and results from a Slice operation.
type SliceReport struct {
	Timestamp time.Time      `json:"timestamp"`
	Offset    int            `json:"offset"`
	Limit     int            `json:"limit"`
	Reverse   bool           `json:"reverse"`
	Total     int            `json:"total_input"`
	Returned  int            `json:"returned"`
	Ports     []scanner.Port `json:"ports"`
}

// BuildSliceReport constructs a SliceReport from input ports and options.
func BuildSliceReport(ports []scanner.Port, opts snapshot.SliceOptions) SliceReport {
	result := snapshot.Slice(ports, opts)
	return SliceReport{
		Timestamp: time.Now().UTC(),
		Offset:    opts.Offset,
		Limit:     opts.Limit,
		Reverse:   opts.Reverse,
		Total:     len(ports),
		Returned:  len(result),
		Ports:     result,
	}
}

// WriteSliceText writes a human-readable slice report to w.
func WriteSliceText(w io.Writer, r SliceReport) error {
	fmt.Fprintf(w, "Slice Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  Input:    %d ports\n", r.Total)
	fmt.Fprintf(w, "  Returned: %d ports\n", r.Returned)
	fmt.Fprintf(w, "  Offset:   %d  Limit: %d  Reverse: %v\n", r.Offset, r.Limit, r.Reverse)
	fmt.Fprintln(w, "  --- Ports ---")
	for _, p := range r.Ports {
		fmt.Fprintf(w, "  %s/%d  pid=%d  process=%s\n", p.Protocol, p.Port, p.PID, p.Process)
	}
	return nil
}

// WriteSliceJSON writes the slice report as JSON to w.
func WriteSliceJSON(w io.Writer, r SliceReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
