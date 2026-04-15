package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// DiffChainReport summarises a sequential diff chain.
type DiffChainReport struct {
	Timestamp   time.Time              `json:"timestamp"`
	EntryCount  int                    `json:"entry_count"`
	TotalAdded  int                    `json:"total_added"`
	TotalRemoved int                   `json:"total_removed"`
	Entries     []DiffChainEntryReport `json:"entries"`
}

// DiffChainEntryReport is a single row in the chain.
type DiffChainEntryReport struct {
	Index     int       `json:"index"`
	Timestamp time.Time `json:"timestamp"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
}

// BuildDiffChainReport converts a chain into a report.
func BuildDiffChainReport(chain []snapshot.ChainEntry) DiffChainReport {
	r := DiffChainReport{
		Timestamp:  time.Now().UTC(),
		EntryCount: len(chain),
	}
	for _, e := range chain {
		added := len(e.Diff.Added)
		removed := len(e.Diff.Removed)
		r.TotalAdded += added
		r.TotalRemoved += removed
		r.Entries = append(r.Entries, DiffChainEntryReport{
			Index:     e.Index,
			Timestamp: e.Snap.Timestamp,
			Added:     added,
			Removed:   removed,
		})
	}
	return r
}

// WriteDiffChainText writes a human-readable diff chain report.
func WriteDiffChainText(w io.Writer, r DiffChainReport) error {
	_, err := fmt.Fprintf(w, "Diff Chain Report — %s\n", r.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Entries: %d | Total Added: %d | Total Removed: %d\n",
		r.EntryCount, r.TotalAdded, r.TotalRemoved)
	if err != nil {
		return err
	}
	for _, e := range r.Entries {
		_, err = fmt.Fprintf(w, "  [%d] %s  +%d -%d\n",
			e.Index, e.Timestamp.Format(time.RFC3339), e.Added, e.Removed)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteDiffChainJSON encodes the report as JSON.
func WriteDiffChainJSON(w io.Writer, r DiffChainReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
