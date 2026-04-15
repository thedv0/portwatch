package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourusername/portwatch/internal/snapshot"
)

// ExpireReport summarises the result of an Expire operation.
type ExpireReport struct {
	Timestamp     time.Time          `json:"timestamp"`
	Total         int                `json:"total"`
	RetainedCount int                `json:"retained_count"`
	ExpiredCount  int                `json:"expired_count"`
	Retained      []snapshot.PortState `json:"retained"`
	Expired       []snapshot.PortState `json:"expired"`
}

// BuildExpireReport constructs an ExpireReport from an ExpireResult.
func BuildExpireReport(res snapshot.ExpireResult) ExpireReport {
	retained := res.Retained
	if retained == nil {
		retained = []snapshot.PortState{}
	}
	expired := res.Expired
	if expired == nil {
		expired = []snapshot.PortState{}
	}
	return ExpireReport{
		Timestamp:     time.Now(),
		Total:         res.Total,
		RetainedCount: len(retained),
		ExpiredCount:  len(expired),
		Retained:      retained,
		Expired:       expired,
	}
}

// WriteExpireText writes a human-readable expire report to w.
func WriteExpireText(w io.Writer, r ExpireReport) error {
	_, err := fmt.Fprintf(w,
		"Expire Report [%s]\n  Total: %d  Retained: %d  Expired: %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.Total, r.RetainedCount, r.ExpiredCount,
	)
	if err != nil {
		return err
	}
	if len(r.Expired) > 0 {
		_, err = fmt.Fprintln(w, "  Expired ports:")
		if err != nil {
			return err
		}
		for _, p := range r.Expired {
			_, err = fmt.Fprintf(w, "    %s/%d (pid=%d)\n", p.Protocol, p.Port, p.PID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteExpireJSON writes an expire report as JSON to w.
func WriteExpireJSON(w io.Writer, r ExpireReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
