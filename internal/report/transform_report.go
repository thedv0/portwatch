package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// TransformReport summarises the result of a Transform operation.
type TransformReport struct {
	Timestamp   time.Time      `json:"timestamp"`
	Total       int            `json:"total"`
	Transformed []scanner.Port `json:"transformed"`
}

// BuildTransformReport constructs a TransformReport from the output of
// snapshot.Transform.
func BuildTransformReport(ports []scanner.Port) TransformReport {
	return TransformReport{
		Timestamp:   time.Now().UTC(),
		Total:       len(ports),
		Transformed: ports,
	}
}

// WriteTransformText writes a human-readable summary to w.
func WriteTransformText(w io.Writer, r TransformReport) error {
	_, err := fmt.Fprintf(w, "Transform Report — %s\n", r.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Total ports: %d\n", r.Total)
	if err != nil {
		return err
	}
	if r.Total == 0 {
		_, err = fmt.Fprintln(w, "  (no ports)")
		return err
	}
	for _, p := range r.Transformed {
		process := p.Process
		if process == "" {
			process = "unknown"
		}
		_, err = fmt.Fprintf(w, "  %-6s %5d  pid=%-6d  process=%s\n",
			p.Protocol, p.Port, p.PID, process)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteTransformJSON serialises r as indented JSON to w.
func WriteTransformJSON(w io.Writer, r TransformReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
