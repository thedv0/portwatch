package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// CorrelateReport summarises the results of a port correlation run.
type CorrelateReport struct {
	Timestamp time.Time                    `json:"timestamp"`
	TotalGroups int                        `json:"total_groups"`
	Groups    []snapshot.CorrelatedGroup   `json:"groups"`
}

// BuildCorrelateReport constructs a CorrelateReport from the given groups.
func BuildCorrelateReport(groups []snapshot.CorrelatedGroup) CorrelateReport {
	return CorrelateReport{
		Timestamp:   time.Now().UTC(),
		TotalGroups: len(groups),
		Groups:      groups,
	}
}

// WriteCorrelateText writes a human-readable correlation report to w.
func WriteCorrelateText(w io.Writer, r CorrelateReport) error {
	_, err := fmt.Fprintf(w, "Correlation Report — %s\n", r.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Total groups: %d\n\n", r.TotalGroups)
	if err != nil {
		return err
	}
	for _, g := range r.Groups {
		_, err = fmt.Fprintf(w, "  [%s] occurrences=%d\n", g.Key, g.Count)
		if err != nil {
			return err
		}
		for _, p := range g.Ports {
			_, err = fmt.Fprintf(w, "    port=%d proto=%s process=%s pid=%d\n",
				p.Port, p.Protocol, p.Process, p.PID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteCorrelateJSON writes the report as JSON to w.
func WriteCorrelateJSON(w io.Writer, r CorrelateReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
