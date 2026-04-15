package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// SplitReport summarises the output of a snapshot.Split call.
type SplitReport struct {
	Timestamp  time.Time              `json:"timestamp"`
	Field      string                 `json:"field"`
	PartCount  int                    `json:"part_count"`
	TotalPorts int                    `json:"total_ports"`
	Parts      []snapshot.SplitResult `json:"parts"`
}

// BuildSplitReport creates a SplitReport from split results.
func BuildSplitReport(results []snapshot.SplitResult, field string) SplitReport {
	total := 0
	for _, r := range results {
		total += len(r.Ports)
	}
	return SplitReport{
		Timestamp:  time.Now().UTC(),
		Field:      field,
		PartCount:  len(results),
		TotalPorts: total,
		Parts:      results,
	}
}

// WriteSplitText writes a human-readable split report to w.
func WriteSplitText(w io.Writer, r SplitReport) error {
	_, err := fmt.Fprintf(w, "Split Report\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "  Timestamp : %s\n", r.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "  Field     : %s\n", r.Field)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "  Parts     : %d\n", r.PartCount)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "  Total     : %d ports\n", r.TotalPorts)
	if err != nil {
		return err
	}
	for _, part := range r.Parts {
		_, err = fmt.Fprintf(w, "  [%d] %d ports\n", part.Index, len(part.Ports))
		if err != nil {
			return err
		}
		for _, p := range part.Ports {
			_, err = fmt.Fprintf(w, "      %-6d %s\t%s\n", p.Port, p.Protocol, p.Process)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteSplitJSON encodes the report as JSON to w.
func WriteSplitJSON(w io.Writer, r SplitReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
