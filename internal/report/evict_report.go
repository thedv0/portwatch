package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// EvictReport summarises the result of an eviction pass.
type EvictReport struct {
	Timestamp  time.Time            `json:"timestamp"`
	Policy     string               `json:"policy"`
	Before     int                  `json:"before"`
	After      int                  `json:"after"`
	Evicted    int                  `json:"evicted"`
	Retained   []snapshot.PortState `json:"retained"`
}

// BuildEvictReport constructs an EvictReport from before/after port slices.
func BuildEvictReport(policy string, before, after []snapshot.PortState) EvictReport {
	return EvictReport{
		Timestamp: time.Now().UTC(),
		Policy:    policy,
		Before:    len(before),
		After:     len(after),
		Evicted:   len(before) - len(after),
		Retained:  after,
	}
}

// WriteEvictText writes a human-readable eviction summary to w.
func WriteEvictText(w io.Writer, r EvictReport) error {
	_, err := fmt.Fprintf(w,
		"Evict Report [%s]\n"+
			"  Policy:  %s\n"+
			"  Before:  %d\n"+
			"  After:   %d\n"+
			"  Evicted: %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.Policy,
		r.Before,
		r.After,
		r.Evicted,
	)
	if err != nil {
		return err
	}
	for _, p := range r.Retained {
		if _, e := fmt.Fprintf(w, "  retained  %s/%d\n", p.Protocol, p.Port); e != nil {
			return e
		}
	}
	return nil
}

// WriteEvictJSON serialises the report as JSON.
func WriteEvictJSON(w io.Writer, r EvictReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
