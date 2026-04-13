package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/username/portwatch/internal/snapshot"
)

// IndexReport summarises the result of an index operation.
type IndexReport struct {
	Timestamp   time.Time              `json:"timestamp"`
	KeyFields   []string               `json:"key_fields"`
	TotalKeys   int                    `json:"total_keys"`
	TotalPorts  int                    `json:"total_ports"`
	Entries     []IndexReportEntry     `json:"entries"`
}

// IndexReportEntry is one row in the report.
type IndexReportEntry struct {
	Key       string `json:"key"`
	PortCount int    `json:"port_count"`
}

// BuildIndexReport constructs an IndexReport from an index map.
func BuildIndexReport(idx map[string]snapshot.IndexEntry, keyFields []string) IndexReport {
	keys := snapshot.SortedKeys(idx)
	entries := make([]IndexReportEntry, 0, len(keys))
	total := 0
	for _, k := range keys {
		e := idx[k]
		entries = append(entries, IndexReportEntry{
			Key:       k,
			PortCount: len(e.Ports),
		})
		total += len(e.Ports)
	}
	return IndexReport{
		Timestamp:  time.Now(),
		KeyFields:  keyFields,
		TotalKeys:  len(keys),
		TotalPorts: total,
		Entries:    entries,
	}
}

// WriteIndexText writes a human-readable index report to w.
func WriteIndexText(w io.Writer, r IndexReport) {
	fmt.Fprintf(w, "Index Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Key fields : %v\n", r.KeyFields)
	fmt.Fprintf(w, "Total keys : %d\n", r.TotalKeys)
	fmt.Fprintf(w, "Total ports: %d\n", r.TotalPorts)
	fmt.Fprintln(w, "---")
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  %-30s %d port(s)\n", e.Key, e.PortCount)
	}
}

// WriteIndexJSON writes the report as JSON to w.
func WriteIndexJSON(w io.Writer, r IndexReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
