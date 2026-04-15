package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/portwatch/internal/snapshot"
)

// SealReport summarises the result of a seal or verify operation.
type SealReport struct {
	Timestamp time.Time `json:"timestamp"`
	Hash      string    `json:"hash"`
	SealedAt  time.Time `json:"sealed_at"`
	PortCount int       `json:"port_count"`
	Verified  bool      `json:"verified"`
	Error     string    `json:"error,omitempty"`
}

// BuildSealReport constructs a SealReport from a SealedSnapshot and an
// optional verification error.
func BuildSealReport(sealed snapshot.SealedSnapshot, verifyErr error) SealReport {
	r := SealReport{
		Timestamp: time.Now(),
		Hash:      sealed.Hash,
		SealedAt:  sealed.SealedAt,
		PortCount: len(sealed.Snapshot.Ports),
		Verified:  verifyErr == nil,
	}
	if verifyErr != nil {
		r.Error = verifyErr.Error()
	}
	return r
}

// WriteSealText writes a human-readable seal report to w.
func WriteSealText(w io.Writer, r SealReport) error {
	fmt.Fprintf(w, "Seal Report\n")
	fmt.Fprintf(w, "-----------\n")
	fmt.Fprintf(w, "Hash      : %s\n", r.Hash)
	fmt.Fprintf(w, "Sealed At : %s\n", r.SealedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Ports     : %d\n", r.PortCount)
	status := "OK"
	if !r.Verified {
		status = fmt.Sprintf("FAILED (%s)", r.Error)
	}
	fmt.Fprintf(w, "Verified  : %s\n", status)
	return nil
}

// WriteSealJSON writes a JSON-encoded seal report to w.
func WriteSealJSON(w io.Writer, r SealReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
