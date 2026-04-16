package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// TruncateReport summarises the result of a truncation pass.
type TruncateReport struct {
	Timestamp   time.Time           `json:"timestamp"`
	InputSnaps  int                 `json:"input_snaps"`
	OutputSnaps int                 `json:"output_snaps"`
	PortsBefore int                 `json:"ports_before"`
	PortsAfter  int                 `json:"ports_after"`
	Options     snapshot.TruncateOptions `json:"options"`
}

// BuildTruncateReport constructs a report from before/after snapshot slices.
func BuildTruncateReport(before, after []snapshot.Snapshot, opts snapshot.TruncateOptions) TruncateReport {
	portsBefore := 0
	for _, s := range before {
		portsBefore += len(s.Ports)
	}
	portsAfter := 0
	for _, s := range after {
		portsAfter += len(s.Ports)
	}
	return TruncateReport{
		Timestamp:   time.Now().UTC(),
		InputSnaps:  len(before),
		OutputSnaps: len(after),
		PortsBefore: portsBefore,
		PortsAfter:  portsAfter,
		Options:     opts,
	}
}

// WriteTruncateText writes a human-readable truncation summary to w.
func WriteTruncateText(w io.Writer, r TruncateReport) error {
	_, err := fmt.Fprintf(w,
		"Truncate Report\n"+
			"  Timestamp:    %s\n"+
			"  Snapshots:    %d -> %d\n"+
			"  Ports:        %d -> %d (removed %d)\n",
		r.Timestamp.Format(time.RFC3339),
		r.InputSnaps, r.OutputSnaps,
		r.PortsBefore, r.PortsAfter, r.PortsBefore-r.PortsAfter,
	)
	return err
}

// WriteTruncateJSON serialises the report as JSON to w.
func WriteTruncateJSON(w io.Writer, r TruncateReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
