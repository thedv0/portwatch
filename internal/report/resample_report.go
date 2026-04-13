package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ResampleReport is the structured output for a resample operation.
type ResampleReport struct {
	GeneratedAt time.Time        `json:"generated_at"`
	BucketSize  string           `json:"bucket_size"`
	Aggregator  string           `json:"aggregator"`
	BucketCount int              `json:"bucket_count"`
	Buckets     []BucketSummary  `json:"buckets"`
}

// BucketSummary holds per-bucket metadata.
type BucketSummary struct {
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	PortCount int       `json:"port_count"`
}

// BuildResampleReport constructs a ResampleReport from resampled buckets.
func BuildResampleReport(buckets []snapshot.Bucket, opts snapshot.ResampleOptions) ResampleReport {
	summaries := make([]BucketSummary, 0, len(buckets))
	for _, b := range buckets {
		summaries = append(summaries, BucketSummary{
			Start:     b.Start,
			End:       b.End,
			PortCount: len(b.Ports),
		})
	}
	return ResampleReport{
		GeneratedAt: time.Now().UTC(),
		BucketSize:  opts.BucketSize.String(),
		Aggregator:  opts.Aggregator,
		BucketCount: len(buckets),
		Buckets:     summaries,
	}
}

// WriteResampleText writes a human-readable resample report to w.
func WriteResampleText(w io.Writer, r ResampleReport) {
	fmt.Fprintf(w, "Resample Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Bucket Size : %s\n", r.BucketSize)
	fmt.Fprintf(w, "Aggregator  : %s\n", r.Aggregator)
	fmt.Fprintf(w, "Buckets     : %d\n\n", r.BucketCount)
	for i, b := range r.Buckets {
		fmt.Fprintf(w, "  [%d] %s → %s  ports=%d\n",
			i+1,
			b.Start.Format(time.RFC3339),
			b.End.Format(time.RFC3339),
			b.PortCount,
		)
	}
}

// WriteResampleJSON writes a JSON-encoded resample report to w.
func WriteResampleJSON(w io.Writer, r ResampleReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
