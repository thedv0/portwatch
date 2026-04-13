package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// PivotReport summarises a pivot operation for reporting.
type PivotReport struct {
	Timestamp time.Time                  `json:"timestamp"`
	Field     snapshot.PivotField        `json:"field"`
	Buckets   map[string]int             `json:"buckets"` // key -> port count
	Keys      []string                   `json:"keys"`
	Total     int                        `json:"total_ports"`
}

// BuildPivotReport converts a PivotResult into a PivotReport.
func BuildPivotReport(res snapshot.PivotResult) PivotReport {
	buckets := make(map[string]int, len(res.Buckets))
	total := 0
	for k, ports := range res.Buckets {
		buckets[k] = len(ports)
		total += len(ports)
	}
	return PivotReport{
		Timestamp: time.Now().UTC(),
		Field:     res.Field,
		Buckets:   buckets,
		Keys:      res.Keys,
		Total:     total,
	}
}

// WritePivotText writes a human-readable pivot report to w.
func WritePivotText(w io.Writer, r PivotReport) error {
	fmt.Fprintf(w, "Pivot Report — field: %s\n", r.Field)
	fmt.Fprintf(w, "Timestamp : %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Total ports: %d\n", r.Total)
	fmt.Fprintln(w, "---")
	for _, k := range r.Keys {
		fmt.Fprintf(w, "  %-30s %d\n", k, r.Buckets[k])
	}
	return nil
}

// WritePivotJSON writes a JSON pivot report to w.
func WritePivotJSON(w io.Writer, r PivotReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
