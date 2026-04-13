package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ThrottleEntry describes a single throttled port event.
type ThrottleEntry struct {
	Key         string    `json:"key"`
	Suppressed  int       `json:"suppressed"`
	LastAllowed time.Time `json:"last_allowed"`
}

// ThrottleReport summarises which port events were throttled.
type ThrottleReport struct {
	Timestamp    time.Time       `json:"timestamp"`
	TotalAllowed int             `json:"total_allowed"`
	TotalBlocked int             `json:"total_blocked"`
	Entries      []ThrottleEntry `json:"entries"`
}

// BuildThrottleReport constructs a ThrottleReport from raw counters.
func BuildThrottleReport(allowed map[string]time.Time, blocked map[string]int) ThrottleReport {
	r := ThrottleReport{
		Timestamp: time.Now().UTC(),
	}
	for key, last := range allowed {
		r.TotalAllowed++
		entry := ThrottleEntry{Key: key, LastAllowed: last}
		if n, ok := blocked[key]; ok {
			entry.Suppressed = n
			r.TotalBlocked += n
		}
		r.Entries = append(r.Entries, entry)
	}
	return r
}

// WriteThrottleText writes a human-readable throttle report to w.
func WriteThrottleText(w io.Writer, r ThrottleReport) error {
	fmt.Fprintf(w, "Throttle Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  Allowed : %d\n", r.TotalAllowed)
	fmt.Fprintf(w, "  Blocked : %d\n", r.TotalBlocked)
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "  (no entries)")
		return nil
	}
	fmt.Fprintln(w, "  Key                     Suppressed  Last Allowed")
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  %-24s %-11d %s\n", e.Key, e.Suppressed, e.LastAllowed.Format(time.RFC3339))
	}
	return nil
}

// WriteThrottleJSON writes the report as JSON to w.
func WriteThrottleJSON(w io.Writer, r ThrottleReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
