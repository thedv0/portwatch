package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// RetainReport summarises the result of a Retain operation.
type RetainReport struct {
	Timestamp   time.Time           `json:"timestamp"`
	InputCount  int                 `json:"input_count"`
	OutputCount int                 `json:"output_count"`
	Dropped     int                 `json:"dropped"`
	Options     snapshot.RetainOptions `json:"options"`
	Snapshots   []snapshot.Snapshot `json:"snapshots"`
}

// BuildRetainReport constructs a RetainReport from the input and output snapshots.
func BuildRetainReport(input, output []snapshot.Snapshot, opts snapshot.RetainOptions) RetainReport {
	return RetainReport{
		Timestamp:   time.Now().UTC(),
		InputCount:  len(input),
		OutputCount: len(output),
		Dropped:     len(input) - len(output),
		Options:     opts,
		Snapshots:   output,
	}
}

// WriteRetainText writes a human-readable summary to w.
func WriteRetainText(w io.Writer, r RetainReport) error {
	_, err := fmt.Fprintf(w,
		"Retain Report\n"+
			"  Timestamp : %s\n"+
			"  Input     : %d snapshots\n"+
			"  Output    : %d snapshots\n"+
			"  Dropped   : %d snapshots\n"+
			"  MaxAge    : %s\n"+
			"  MinCount  : %d\n"+
			"  MaxCount  : %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.InputCount,
		r.OutputCount,
		r.Dropped,
		r.Options.MaxAge,
		r.Options.MinCount,
		r.Options.MaxCount,
	)
	return err
}

// WriteRetainJSON writes the report as JSON to w.
func WriteRetainJSON(w io.Writer, r RetainReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
