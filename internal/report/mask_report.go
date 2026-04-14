package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// MaskReport holds the result of a masking operation.
type MaskReport struct {
	Timestamp time.Time            `json:"timestamp"`
	Total     int                  `json:"total"`
	Options   snapshot.MaskOptions `json:"options"`
	Ports     []snapshot.PortEntry `json:"ports"`
}

// BuildMaskReport constructs a MaskReport from masked port entries and the
// options that were applied.
func BuildMaskReport(ports []snapshot.PortEntry, opts snapshot.MaskOptions) MaskReport {
	return MaskReport{
		Timestamp: time.Now().UTC(),
		Total:     len(ports),
		Options:   opts,
		Ports:     ports,
	}
}

// WriteMaskText writes a human-readable mask report to w.
func WriteMaskText(w io.Writer, r MaskReport) error {
	fmt.Fprintf(w, "Mask Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Total entries : %d\n", r.Total)
	fmt.Fprintf(w, "MaskProcess   : %v\n", r.Options.MaskProcess)
	fmt.Fprintf(w, "MaskPID       : %v\n", r.Options.MaskPID)
	fmt.Fprintf(w, "MaskPort      : %v\n", r.Options.MaskPort)
	fmt.Fprintln(w, "---")
	for _, p := range r.Ports {
		fmt.Fprintf(w, "  port=%-6d pid=%-6d proto=%-5s process=%s\n",
			p.Port, p.PID, p.Protocol, p.Process)
	}
	return nil
}

// WriteMaskJSON writes the report as JSON to w.
func WriteMaskJSON(w io.Writer, r MaskReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
