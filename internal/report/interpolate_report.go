package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// InterpolateReport summarises the result of a snapshot interpolation pass.
type InterpolateReport struct {
	Timestamp      time.Time `json:"timestamp"`
	TotalSnapshots int       `json:"total_snapshots"`
	Filled         int       `json:"filled"`
	Method         string    `json:"method"`
	Step           string    `json:"step"`
}

// BuildInterpolateReport constructs a report from an InterpolateResult.
func BuildInterpolateReport(res snapshot.InterpolateResult, method, step string) InterpolateReport {
	return InterpolateReport{
		Timestamp:      res.Timestamp,
		TotalSnapshots: len(res.Snapshots),
		Filled:         res.Filled,
		Method:         method,
		Step:           step,
	}
}

// WriteInterpolateText writes a human-readable interpolation report to w.
func WriteInterpolateText(w io.Writer, r InterpolateReport) error {
	_, err := fmt.Fprintf(w,
		"Interpolate Report\n"+
			"  Timestamp : %s\n"+
			"  Method    : %s\n"+
			"  Step      : %s\n"+
			"  Total     : %d\n"+
			"  Filled    : %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.Method,
		r.Step,
		r.TotalSnapshots,
		r.Filled,
	)
	return err
}

// WriteInterpolateJSON writes the report as JSON to w.
func WriteInterpolateJSON(w io.Writer, r InterpolateReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
