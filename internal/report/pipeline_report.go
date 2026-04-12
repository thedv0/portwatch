package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/netwatch/portwatch/internal/snapshot"
)

// PipelineReport summarises the output of a pipeline run.
type PipelineReport struct {
	Timestamp    time.Time `json:"timestamp"`
	TotalInput   int       `json:"total_input"`
	TotalOutput  int       `json:"total_output"`
	Duplicates   int       `json:"duplicates_removed"`
	Issues       int       `json:"validation_issues"`
	Classified   int       `json:"classified_ports"`
}

// BuildPipelineReport constructs a PipelineReport from a result and the
// original input length.
func BuildPipelineReport(inputLen int, result snapshot.PipelineResult) PipelineReport {
	issues := 0
	if result.Validation != nil {
		issues = len(result.Validation.Issues)
	}
	return PipelineReport{
		Timestamp:   time.Now().UTC(),
		TotalInput:  inputLen,
		TotalOutput: len(result.Ports),
		Duplicates:  inputLen - len(result.Ports),
		Issues:      issues,
		Classified:  len(result.Classified),
	}
}

// WritePipelineText writes a human-readable summary to w.
func WritePipelineText(w io.Writer, r PipelineReport) error {
	_, err := fmt.Fprintf(w,
		"Pipeline Report [%s]\n  Input:      %d\n  Output:     %d\n  Duplicates: %d\n  Issues:     %d\n  Classified: %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.TotalInput, r.TotalOutput, r.Duplicates, r.Issues, r.Classified,
	)
	return err
}

// WritePipelineJSON writes the report as JSON to w.
func WritePipelineJSON(w io.Writer, r PipelineReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
