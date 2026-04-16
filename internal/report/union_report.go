package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// UnionReport holds the result of a union operation.
type UnionReport struct {
	Timestamp   time.Time      `json:"timestamp"`
	InputSnaps  int            `json:"input_snaps"`
	TotalPorts  int            `json:"total_ports"`
	DedupEnabled bool          `json:"dedup_enabled"`
	KeyFields   []string       `json:"key_fields"`
	Ports       []scanner.Port `json:"ports"`
}

// BuildUnionReport constructs a UnionReport from the given inputs.
func BuildUnionReport(snaps []snapshot.Snapshot, ports []scanner.Port, opts snapshot.UnionOptions) UnionReport {
	return UnionReport{
		Timestamp:    time.Now().UTC(),
		InputSnaps:   len(snaps),
		TotalPorts:   len(ports),
		DedupEnabled: opts.Dedup,
		KeyFields:    opts.KeyFields,
		Ports:        ports,
	}
}

// WriteUnionText writes a human-readable union report to w.
func WriteUnionText(w io.Writer, r UnionReport) error {
	fmt.Fprintf(w, "Union Report\n")
	fmt.Fprintf(w, "============\n")
	fmt.Fprintf(w, "Timestamp:     %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Input Snaps:   %d\n", r.InputSnaps)
	fmt.Fprintf(w, "Total Ports:   %d\n", r.TotalPorts)
	fmt.Fprintf(w, "Dedup Enabled: %v\n", r.DedupEnabled)
	fmt.Fprintf(w, "Key Fields:    %v\n", r.KeyFields)
	fmt.Fprintf(w, "\nPorts:\n")
	for _, p := range r.Ports {
		fmt.Fprintf(w, "  %s/%d  pid=%-6d  %s\n", p.Protocol, p.Port, p.PID, p.Process)
	}
	return nil
}

// WriteUnionJSON writes a JSON-encoded union report to w.
func WriteUnionJSON(w io.Writer, r UnionReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
