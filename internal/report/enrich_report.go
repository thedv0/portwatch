package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// EnrichReport holds enriched port data for reporting.
type EnrichReport struct {
	Timestamp time.Time      `json:"timestamp"`
	Total     int            `json:"total"`
	Ports     []scanner.Port `json:"ports"`
}

// BuildEnrichReport constructs an EnrichReport from enriched ports.
func BuildEnrichReport(ports []scanner.Port) EnrichReport {
	return EnrichReport{
		Timestamp: time.Now().UTC(),
		Total:     len(ports),
		Ports:     ports,
	}
}

// WriteEnrichText writes a human-readable enrichment report to w.
func WriteEnrichText(w io.Writer, r EnrichReport) error {
	_, err := fmt.Fprintf(w, "Enrichment Report — %s\n", r.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Total ports: %d\n\n", r.Total)
	if err != nil {
		return err
	}
	for _, p := range r.Ports {
		process := p.Process
		if process == "" {
			process = "(unknown)"
		}
		_, err = fmt.Fprintf(w, "  %-6s %5d  pid=%-6d  process=%s\n",
			p.Protocol, p.Port, p.PID, process)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteEnrichJSON writes the report as JSON to w.
func WriteEnrichJSON(w io.Writer, r EnrichReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
