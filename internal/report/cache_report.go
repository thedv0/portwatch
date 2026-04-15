package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// CacheReport summarises the state of a snapshot cache.
type CacheReport struct {
	Timestamp   time.Time `json:"timestamp"`
	LiveEntries int       `json:"live_entries"`
	MaxEntries  int       `json:"max_entries"`
	TTLSeconds  float64   `json:"ttl_seconds"`
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
}

// BuildCacheReport constructs a CacheReport from the provided stats.
func BuildCacheReport(live, max int, ttl time.Duration, hits, misses int64) CacheReport {
	return CacheReport{
		Timestamp:   time.Now(),
		LiveEntries: live,
		MaxEntries:  max,
		TTLSeconds:  ttl.Seconds(),
		Hits:        hits,
		Misses:      misses,
	}
}

// WriteCacheText writes a human-readable cache report to w.
func WriteCacheText(w io.Writer, r CacheReport) error {
	_, err := fmt.Fprintf(w,
		"Cache Report [%s]\n"+
			"  Live Entries : %d / %d\n"+
			"  TTL          : %.0fs\n"+
			"  Hits         : %d\n"+
			"  Misses       : %d\n",
		r.Timestamp.Format(time.RFC3339),
		r.LiveEntries, r.MaxEntries,
		r.TTLSeconds,
		r.Hits,
		r.Misses,
	)
	return err
}

// WriteCacheJSON writes the cache report as JSON to w.
func WriteCacheJSON(w io.Writer, r CacheReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
