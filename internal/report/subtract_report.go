package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// SubtractReport summarises the result of a Subtract operation.
type SubtractReport struct {
	Timestamp   time.Time      `json:"timestamp"`
	LeftTotal   int            `json:"left_total"`
	RightTotal  int            `json:"right_total"`
	ResultTotal int            `json:"result_total"`
	Removed     int            `json:"removed"`
	Ports       []scanner.Port `json:"ports"`
}

// BuildSubtractReport constructs a SubtractReport from the input and result slices.
func BuildSubtractReport(left, right, result []scanner.Port) SubtractReport {
	return SubtractReport{
		Timestamp:   time.Now().UTC(),
		LeftTotal:   len(left),
		RightTotal:  len(right),
		ResultTotal: len(result),
		Removed:     len(left) - len(result),
		Ports:       result,
	}
}

// WriteSubtractText writes a human-readable subtract report to w.
func WriteSubtractText(w io.Writer, r SubtractReport) error {
	_, err := fmt.Fprintf(w,
		"Subtract Report [%s]\n  Left: %d  Right (excluded): %d  Removed: %d  Result: %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.LeftTotal, r.RightTotal, r.Removed, r.ResultTotal,
	)
	if err != nil {
		return err
	}
	for _, p := range r.Ports {
		if _, err := fmt.Fprintf(w, "  %s:%d pid=%d process=%s\n",
			p.Protocol, p.Port, p.PID, p.Process); err != nil {
			return err
		}
	}
	return nil
}

// WriteSubtractJSON encodes the report as JSON into w.
func WriteSubtractJSON(w io.Writer, r SubtractReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
