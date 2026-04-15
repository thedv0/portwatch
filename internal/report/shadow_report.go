package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/username/portwatch/internal/snapshot"
)

// ShadowReport is the serialisable form of a shadow comparison result.
type ShadowReport struct {
	Timestamp   time.Time                    `json:"timestamp"`
	Diverged    bool                         `json:"diverged"`
	PrimaryLen  int                          `json:"primary_len"`
	ShadowLen   int                          `json:"shadow_len"`
	Divergences []snapshot.ShadowDivergence  `json:"divergences"`
}

// BuildShadowReport converts a ShadowResult into a ShadowReport.
func BuildShadowReport(res snapshot.ShadowResult) ShadowReport {
	divs := res.Divergences
	if divs == nil {
		divs = []snapshot.ShadowDivergence{}
	}
	return ShadowReport{
		Timestamp:   res.Timestamp,
		Diverged:    res.Diverged,
		PrimaryLen:  res.PrimaryLen,
		ShadowLen:   res.ShadowLen,
		Divergences: divs,
	}
}

// WriteShadowText writes a human-readable shadow report to w.
func WriteShadowText(w io.Writer, r ShadowReport) error {
	fmt.Fprintf(w, "Shadow Comparison Report\n")
	fmt.Fprintf(w, "Timestamp : %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Diverged  : %v\n", r.Diverged)
	fmt.Fprintf(w, "Primary   : %d ports\n", r.PrimaryLen)
	fmt.Fprintf(w, "Shadow    : %d ports\n", r.ShadowLen)
	fmt.Fprintf(w, "Divergences (%d):\n", len(r.Divergences))
	if len(r.Divergences) == 0 {
		fmt.Fprintf(w, "  (none)\n")
		return nil
	}
	for _, d := range r.Divergences {
		fmt.Fprintf(w, "  [%s:%d] %s\n", d.Protocol, d.Port, d.Reason)
	}
	return nil
}

// WriteShadowJSON writes a JSON-encoded shadow report to w.
func WriteShadowJSON(w io.Writer, r ShadowReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
