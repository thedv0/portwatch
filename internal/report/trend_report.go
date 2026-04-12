package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// TrendEntry pairs a metric name with its computed trend.
type TrendEntry struct {
	Metric    string                  `json:"metric"`
	Direction snapshot.TrendDirection `json:"direction"`
	Slope     float64                 `json:"slope"`
	Points    int                     `json:"points"`
}

// TrendReport is the top-level structure for a trend analysis report.
type TrendReport struct {
	GeneratedAt time.Time     `json:"generated_at"`
	Entries     []TrendEntry  `json:"entries"`
}

// BuildTrendReport constructs a TrendReport from a map of metric name to
// trend points, using the provided options.
func BuildTrendReport(
	metrics map[string][]snapshot.TrendPoint,
	opts snapshot.TrendOptions,
) TrendReport {
	report := TrendReport{GeneratedAt: time.Now()}
	for name, pts := range metrics {
		res := snapshot.AnalyzeTrend(pts, opts)
		report.Entries = append(report.Entries, TrendEntry{
			Metric:    name,
			Direction: res.Direction,
			Slope:     res.Slope,
			Points:    len(res.Points),
		})
	}
	return report
}

// WriteTrendText writes a human-readable trend report to w.
func WriteTrendText(w io.Writer, r TrendReport) error {
	fmt.Fprintf(w, "Trend Report — %s\n\n", r.GeneratedAt.Format(time.RFC3339))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tDIRECTION\tSLOPE\tPOINTS")
	for _, e := range r.Entries {
		fmt.Fprintf(tw, "%s\t%s\t%.4f\t%d\n", e.Metric, e.Direction, e.Slope, e.Points)
	}
	return tw.Flush()
}

// WriteTrendJSON writes the trend report as JSON to w.
func WriteTrendJSON(w io.Writer, r TrendReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
