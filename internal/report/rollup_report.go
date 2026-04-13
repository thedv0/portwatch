package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// RollupReport is the structured output of a rollup operation.
type RollupReport struct {
	Timestamp     time.Time        `json:"timestamp"`
	SnapshotCount int              `json:"snapshot_count"`
	TotalPorts    int              `json:"total_ports"`
	Dropped       int              `json:"dropped"`
	Ports         []RollupPortRow  `json:"ports"`
}

// RollupPortRow is a single port entry in the rollup report.
type RollupPortRow struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Process  string `json:"process"`
	PID      int    `json:"pid"`
}

// BuildRollupReport constructs a RollupReport from a RollupResult.
func BuildRollupReport(res snapshot.RollupResult) RollupReport {
	rows := make([]RollupPortRow, 0, len(res.Ports))
	for _, p := range res.Ports {
		rows = append(rows, RollupPortRow{
			Port:     p.Port,
			Protocol: p.Protocol,
			Process:  p.Process,
			PID:      p.PID,
		})
	}
	return RollupReport{
		Timestamp:     res.Timestamp,
		SnapshotCount: res.SnapshotCount,
		TotalPorts:    len(res.Ports),
		Dropped:       res.Dropped,
		Ports:         rows,
	}
}

// WriteRollupText writes a human-readable rollup report to w.
func WriteRollupText(w io.Writer, r RollupReport) {
	fmt.Fprintf(w, "Rollup Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Snapshots included : %d\n", r.SnapshotCount)
	fmt.Fprintf(w, "Total ports        : %d\n", r.TotalPorts)
	fmt.Fprintf(w, "Dropped (threshold): %d\n", r.Dropped)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%-8s %-10s %-20s %s\n", "PORT", "PROTOCOL", "PROCESS", "PID")
	for _, p := range r.Ports {
		fmt.Fprintf(w, "%-8d %-10s %-20s %d\n", p.Port, p.Protocol, p.Process, p.PID)
	}
}

// WriteRollupJSON writes the rollup report as JSON to w.
func WriteRollupJSON(w io.Writer, r RollupReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
