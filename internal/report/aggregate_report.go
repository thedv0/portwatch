package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// AggregateReport wraps AggregateStats for rendering.
type AggregateReport struct {
	GeneratedAt time.Time                 `json:"generated_at"`
	Stats       snapshot.AggregateStats   `json:"stats"`
}

// WriteAggregateText renders an aggregate report as human-readable text.
func WriteAggregateText(w io.Writer, r AggregateReport) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Aggregate Report\t\n")
	fmt.Fprintf(tw, "Generated:\t%s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(tw, "Period:\t%s — %s\n",
		r.Stats.From.Format(time.RFC3339),
		r.Stats.To.Format(time.RFC3339))
	fmt.Fprintf(tw, "Samples:\t%d\n", r.Stats.SampleCount)
	fmt.Fprintf(tw, "Avg Open Ports:\t%.2f\n", r.Stats.AvgOpen)
	fmt.Fprintf(tw, "Max Open Ports:\t%d\n", r.Stats.MaxOpen)
	fmt.Fprintf(tw, "Min Open Ports:\t%d\n", r.Stats.MinOpen)

	if len(r.Stats.TopPorts) > 0 {
		fmt.Fprintf(tw, "\nTop Ports:\t\n")
		for _, p := range r.Stats.TopPorts {
			fmt.Fprintf(tw, "  :%d\t(seen %d times)\n", p.Port, p.Count)
		}
	}

	if len(r.Stats.TopProcesses) > 0 {
		fmt.Fprintf(tw, "\nTop Processes:\t\n")
		for _, p := range r.Stats.TopProcesses {
			fmt.Fprintf(tw, "  %s\t(seen %d times)\n", p.Process, p.Count)
		}
	}

	return tw.Flush()
}

// WriteAggregateJSON renders an aggregate report as JSON.
func WriteAggregateJSON(w io.Writer, r AggregateReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// BuildAggregateReport constructs an AggregateReport from a slice of snapshots.
func BuildAggregateReport(snaps []snapshot.Snapshot) AggregateReport {
	return AggregateReport{
		GeneratedAt: time.Now().UTC(),
		Stats:       snapshot.Aggregate(snaps),
	}
}
