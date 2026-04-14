package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/example/portwatch/internal/snapshot"
)

// CapReport summarises the result of a Cap operation.
type CapReport struct {
	Timestamp   time.Time            `json:"timestamp"`
	Original    int                  `json:"original"`
	Retained    int                  `json:"retained"`
	Dropped     int                  `json:"dropped"`
	MaxPorts    int                  `json:"max_ports"`
	SortedFirst bool                 `json:"sorted_first"`
	Ports       []snapshot.PortState `json:"ports"`
}

// BuildCapReport constructs a CapReport from before/after port slices and opts.
func BuildCapReport(before, after []snapshot.PortState, opts snapshot.CapOptions) CapReport {
	return CapReport{
		Timestamp:   time.Now().UTC(),
		Original:    len(before),
		Retained:    len(after),
		Dropped:     len(before) - len(after),
		MaxPorts:    opts.MaxPorts,
		SortedFirst: opts.SortByPort,
		Ports:       after,
	}
}

// WriteCapText writes a human-readable cap report to w.
func WriteCapText(w io.Writer, r CapReport) error {
	_, err := fmt.Fprintf(w,
		"Cap Report — %s\n"+
			"  Original : %d\n"+
			"  Retained : %d\n"+
			"  Dropped  : %d\n"+
			"  Max Ports: %d\n"+
			"  Sorted   : %v\n",
		r.Timestamp.Format(time.RFC3339),
		r.Original,
		r.Retained,
		r.Dropped,
		r.MaxPorts,
		r.SortedFirst,
	)
	if err != nil {
		return err
	}
	for _, p := range r.Ports {
		if _, err := fmt.Fprintf(w, "  %s/%d (pid %d) %s\n", p.Protocol, p.Port, p.PID, p.Process); err != nil {
			return err
		}
	}
	return nil
}

// WriteCapJSON encodes the report as JSON to w.
func WriteCapJSON(w io.Writer, r CapReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
