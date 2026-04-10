package snapshot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ExportFormat specifies the output format for snapshot exports.
type ExportFormat string

const (
	FormatCSV  ExportFormat = "csv"
	FormatJSON ExportFormat = "json"
)

// ExportRecord represents a single port entry in an export.
type ExportRecord struct {
	Timestamp time.Time         `json:"timestamp"`
	Port      scanner.PortState `json:"port"`
}

// Exporter writes snapshot data to an io.Writer in the requested format.
type Exporter struct {
	w      io.Writer
	format ExportFormat
}

// NewExporter creates an Exporter that writes to w using the given format.
// If format is empty, FormatJSON is used.
func NewExporter(w io.Writer, format ExportFormat) *Exporter {
	if format == "" {
		format = FormatJSON
	}
	return &Exporter{w: w, format: format}
}

// Write serialises the snapshot's ports to the configured format.
func (e *Exporter) Write(snap Snapshot) error {
	switch e.format {
	case FormatCSV:
		return e.writeCSV(snap)
	case FormatJSON:
		return e.writeJSON(snap)
	default:
		return fmt.Errorf("unsupported export format: %q", e.format)
	}
}

func (e *Exporter) writeJSON(snap Snapshot) error {
	records := make([]ExportRecord, len(snap.Ports))
	for i, p := range snap.Ports {
		records[i] = ExportRecord{Timestamp: snap.Timestamp, Port: p}
	}
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func (e *Exporter) writeCSV(snap Snapshot) error {
	w := csv.NewWriter(e.w)
	if err := w.Write([]string{"timestamp", "protocol", "port", "pid", "process"}); err != nil {
		return err
	}
	ts := snap.Timestamp.UTC().Format(time.RFC3339)
	for _, p := range snap.Ports {
		row := []string{
			ts,
			p.Protocol,
			strconv.Itoa(p.Port),
			strconv.Itoa(p.PID),
			p.Process,
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
