package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ReplayReport summarises a replay sequence for presentation.
type ReplayReport struct {
	GeneratedAt time.Time             `json:"generated_at"`
	FrameCount  int                   `json:"frame_count"`
	Frames      []ReplayFrameSummary  `json:"frames"`
}

// ReplayFrameSummary is the reportable view of a single replay frame.
type ReplayFrameSummary struct {
	Index     int       `json:"index"`
	Timestamp time.Time `json:"timestamp"`
	TotalOpen int       `json:"total_open"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
}

// BuildReplayReport converts replay frames into a ReplayReport.
func BuildReplayReport(frames []snapshot.ReplayFrame) ReplayReport {
	summaries := make([]ReplayFrameSummary, 0, len(frames))
	for _, f := range frames {
		summaries = append(summaries, ReplayFrameSummary{
			Index:     f.Index,
			Timestamp: f.Timestamp,
			TotalOpen: len(f.Ports),
			Added:     len(f.Diff.Added),
			Removed:   len(f.Diff.Removed),
		})
	}
	return ReplayReport{
		GeneratedAt: time.Now().UTC(),
		FrameCount:  len(frames),
		Frames:      summaries,
	}
}

// WriteReplayText writes a human-readable replay report to w.
func WriteReplayText(w io.Writer, r ReplayReport) error {
	_, err := fmt.Fprintf(w, "Replay Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "Frames: %d\n\n", r.FrameCount)
	if err != nil {
		return err
	}
	for _, f := range r.Frames {
		_, err = fmt.Fprintf(w, "  [%d] %s  open=%-4d +%-3d -%-3d\n",
			f.Index, f.Timestamp.Format(time.RFC3339), f.TotalOpen, f.Added, f.Removed)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteReplayJSON writes a JSON-encoded replay report to w.
func WriteReplayJSON(w io.Writer, r ReplayReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
