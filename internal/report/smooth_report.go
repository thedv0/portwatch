package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// SmoothReport holds the result of an EMA smoothing pass.
type SmoothReport struct {
	Timestamp time.Time               `json:"timestamp"`
	Alpha     float64                 `json:"alpha"`
	Points    []snapshot.SmoothedPoint `json:"points"`
}

// BuildSmoothReport constructs a SmoothReport from the given points and alpha.
func BuildSmoothReport(pts []snapshot.SmoothedPoint, alpha float64) SmoothReport {
	return SmoothReport{
		Timestamp: time.Now().UTC(),
		Alpha:     alpha,
		Points:    pts,
	}
}

// WriteSmoothText writes a human-readable smooth report to w.
func WriteSmoothText(w io.Writer, r SmoothReport) error {
	fmt.Fprintf(w, "Smooth Report — alpha=%.2f  points=%d\n", r.Alpha, len(r.Points))
	fmt.Fprintf(w, "%-30s  %8s  %10s\n", "Timestamp", "Raw", "Smoothed")
	for _, p := range r.Points {
		fmt.Fprintf(w, "%-30s  %8d  %10.3f\n",
			p.Timestamp.UTC().Format(time.RFC3339),
			p.RawCount,
			p.Smoothed,
		)
	}
	return nil
}

// WriteSmoothJSON serialises the report as JSON to w.
func WriteSmoothJSON(w io.Writer, r SmoothReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
