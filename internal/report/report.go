package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a summary of a port scan cycle.
type Report struct {
	Timestamp time.Time          `json:"timestamp"`
	Added     []snapshot.Port    `json:"added"`
	Removed   []snapshot.Port    `json:"removed"`
	Total     int                `json:"total_open"`
}

// Writer writes reports to an output destination.
type Writer struct {
	out    io.Writer
	format Format
}

// NewWriter creates a Writer. If out is nil, os.Stdout is used.
func NewWriter(out io.Writer, format Format) *Writer {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Writer{out: out, format: format}
}

// Write outputs the report in the configured format.
func (w *Writer) Write(r Report) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(r)
	default:
		return w.writeText(r)
	}
}

func (w *Writer) writeText(r Report) error {
	fmt.Fprintf(w.out, "[%s] Port scan report — total open: %d\n",
		r.Timestamp.Format(time.RFC3339), r.Total)
	for _, p := range r.Added {
		fmt.Fprintf(w.out, "  + %s/%d (pid %d)\n", p.Protocol, p.Port, p.PID)
	}
	for _, p := range r.Removed {
		fmt.Fprintf(w.out, "  - %s/%d (pid %d)\n", p.Protocol, p.Port, p.PID)
	}
	return nil
}

func (w *Writer) writeJSON(r Report) error {
	enc := json.NewEncoder(w.out)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
