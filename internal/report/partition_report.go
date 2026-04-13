package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// PartitionReport summarises the result of a snapshot partition operation.
type PartitionReport struct {
	Timestamp   time.Time          `json:"timestamp"`
	BucketCount int                `json:"bucket_count"`
	TotalPorts  int                `json:"total_ports"`
	Buckets     []BucketSummary    `json:"buckets"`
}

// BucketSummary holds per-bucket metadata for reporting.
type BucketSummary struct {
	Key       string    `json:"key"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	PortCount int       `json:"port_count"`
}

// BuildPartitionReport converts raw buckets into a PartitionReport.
func BuildPartitionReport(buckets []snapshot.Bucket) PartitionReport {
	rep := PartitionReport{
		Timestamp:   time.Now().UTC(),
		BucketCount: len(buckets),
	}
	for _, b := range buckets {
		rep.TotalPorts += len(b.Ports)
		rep.Buckets = append(rep.Buckets, BucketSummary{
			Key:       b.Key,
			Start:     b.Start,
			End:       b.End,
			PortCount: len(b.Ports),
		})
	}
	return rep
}

// WritePartitionText writes a human-readable partition report to w.
func WritePartitionText(w io.Writer, rep PartitionReport) error {
	fmt.Fprintf(w, "Partition Report — %s\n", rep.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Buckets : %d\n", rep.BucketCount)
	fmt.Fprintf(w, "Total Ports: %d\n\n", rep.TotalPorts)
	for _, b := range rep.Buckets {
		if b.Key == "unknown" {
			fmt.Fprintf(w, "  [unknown]  ports=%d\n", b.PortCount)
		} else {
			fmt.Fprintf(w, "  %s → %s  ports=%d\n",
				b.Start.Format(time.RFC3339),
				b.End.Format(time.RFC3339),
				b.PortCount,
			)
		}
	}
	return nil
}

// WritePartitionJSON encodes the report as JSON to w.
func WritePartitionJSON(w io.Writer, rep PartitionReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}
