package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ProjectReport is the serialisable output of a projection run.
type ProjectReport struct {
	GeneratedAt time.Time                  `json:"generated_at"`
	BaseCount   float64                    `json:"base_count"`
	Slope       float64                    `json:"slope"`
	Steps       int                        `json:"steps"`
	Points      []snapshot.ProjectedPoint  `json:"points"`
}

// BuildProjectReport converts a ProjectResult into a ProjectReport.
func BuildProjectReport(res snapshot.ProjectResult) ProjectReport {
	return ProjectReport{
		GeneratedAt: res.GeneratedAt,
		BaseCount:   res.BaseCount,
		Slope:       res.Slope,
		Steps:       len(res.Points),
		Points:      res.Points,
	}
}

// WriteProjectText writes a human-readable projection report to w.
func WriteProjectText(w io.Writer, r ProjectReport) error {
	fmt.Fprintf(w, "Projection Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Base open ports : %.0f\n", r.BaseCount)
	fmt.Fprintf(w, "Trend slope     : %+.4f ports/step\n", r.Slope)
	fmt.Fprintf(w, "Steps projected : %d\n", r.Steps)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%-30s  %s\n", "Time", "Projected Count")
	fmt.Fprintf(w, "%-30s  %s\n", "----", "---------------")
	for _, p := range r.Points {
		fmt.Fprintf(w, "%-30s  %.2f\n", p.At.Format(time.RFC3339), p.Count)
	}
	return nil
}

// WriteProjectJSON writes a JSON-encoded projection report to w.
func WriteProjectJSON(w io.Writer, r ProjectReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
