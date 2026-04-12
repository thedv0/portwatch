package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ValidateReport is the top-level structure for a validation report.
type ValidateReport struct {
	Timestamp time.Time                   `json:"timestamp"`
	Total     int                         `json:"total_ports"`
	ErrorCount   int                      `json:"error_count"`
	WarningCount int                      `json:"warning_count"`
	Issues    []snapshot.ValidationIssue  `json:"issues"`
}

// BuildValidateReport constructs a ValidateReport from ports and options.
func BuildValidateReport(ports []snapshot.PortEntry, opts snapshot.ValidateOptions) (*ValidateReport, error) {
	res, err := snapshot.Validate(ports, opts)
	if err != nil {
		return nil, err
	}
	r := &ValidateReport{
		Timestamp: time.Now().UTC(),
		Total:     len(ports),
		Issues:    res.Issues,
	}
	for _, issue := range res.Issues {
		switch issue.Level {
		case snapshot.LevelError:
			r.ErrorCount++
		case snapshot.LevelWarning:
			r.WarningCount++
		}
	}
	return r, nil
}

// WriteValidateText writes a human-readable validation report to w.
func WriteValidateText(w io.Writer, r *ValidateReport) {
	fmt.Fprintf(w, "Validation Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Total ports: %d | Errors: %d | Warnings: %d\n", r.Total, r.ErrorCount, r.WarningCount)
	if len(r.Issues) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return
	}
	for _, issue := range r.Issues {
		level := "WARN"
		if issue.Level == snapshot.LevelError {
			level = "ERROR"
		} else if issue.Level == snapshot.LevelInfo {
			level = "INFO"
		}
		fmt.Fprintf(w, "  [%s] port=%d process=%q — %s\n", level, issue.Port, issue.Process, issue.Message)
	}
}

// WriteValidateJSON writes the validation report as JSON to w.
func WriteValidateJSON(w io.Writer, r *ValidateReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
