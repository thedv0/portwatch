package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// PatchReport summarises the result of a patch operation.
type PatchReport struct {
	Timestamp time.Time              `json:"timestamp"`
	Applied   int                    `json:"applied"`
	Skipped   int                    `json:"skipped"`
	Total     int                    `json:"total"`
	Ports     []snapshot.PortState   `json:"ports"`
}

// BuildPatchReport constructs a PatchReport from a PatchResult.
func BuildPatchReport(res snapshot.PatchResult) PatchReport {
	return PatchReport{
		Timestamp: res.Timestamp,
		Applied:   res.Applied,
		Skipped:   res.Skipped,
		Total:     len(res.Ports),
		Ports:     res.Ports,
	}
}

// WritePatchText writes a human-readable patch report to w.
func WritePatchText(w io.Writer, r PatchReport) error {
	_, err := fmt.Fprintf(w, "Patch Report — %s\n", r.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Applied : %d\n", r.Applied)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Skipped : %d\n", r.Skipped)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Total   : %d\n", r.Total)
	if err != nil {
		return err
	}
	for _, p := range r.Ports {
		_, err = fmt.Fprintf(w, "  %-6s %5d  pid=%-6d  %s\n", p.Protocol, p.Port, p.PID, p.Process)
		if err != nil {
			return err
		}
	}
	return nil
}

// WritePatchJSON writes the report as JSON to w.
func WritePatchJSON(w io.Writer, r PatchReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
