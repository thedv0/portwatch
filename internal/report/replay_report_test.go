package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeReplayFrames() []snapshot.ReplayFrame {
	now := time.Now()
	return []snapshot.ReplayFrame{
		{
			Index:     0,
			Timestamp: now,
			Ports:     []snapshot.PortState{{Port: 80, Protocol: "tcp"}},
			Diff:      snapshot.DiffResult{Added: []snapshot.PortState{{Port: 80, Protocol: "tcp"}}},
		},
		{
			Index:     1,
			Timestamp: now.Add(time.Minute),
			Ports:     []snapshot.PortState{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}},
			Diff:      snapshot.DiffResult{Added: []snapshot.PortState{{Port: 443, Protocol: "tcp"}}},
		},
	}
}

func TestBuildReplayReport_FrameCount(t *testing.T) {
	frames := makeReplayFrames()
	r := BuildReplayReport(frames)
	if r.FrameCount != 2 {
		t.Errorf("expected FrameCount=2, got %d", r.FrameCount)
	}
}

func TestBuildReplayReport_SummaryFields(t *testing.T) {
	frames := makeReplayFrames()
	r := BuildReplayReport(frames)
	if r.Frames[0].TotalOpen != 1 {
		t.Errorf("frame 0 TotalOpen want 1, got %d", r.Frames[0].TotalOpen)
	}
	if r.Frames[1].Added != 1 {
		t.Errorf("frame 1 Added want 1, got %d", r.Frames[1].Added)
	}
}

func TestWriteReplayText_ContainsHeaders(t *testing.T) {
	r := BuildReplayReport(makeReplayFrames())
	var buf bytes.Buffer
	if err := WriteReplayText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Replay Report") {
		t.Error("output missing 'Replay Report' header")
	}
	if !strings.Contains(out, "Frames:") {
		t.Error("output missing 'Frames:' line")
	}
}

func TestWriteReplayJSON_ValidJSON(t *testing.T) {
	r := BuildReplayReport(makeReplayFrames())
	var buf bytes.Buffer
	if err := WriteReplayJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out ReplayReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.FrameCount != 2 {
		t.Errorf("decoded FrameCount want 2, got %d", out.FrameCount)
	}
}

func TestBuildReplayReport_EmptyInput(t *testing.T) {
	r := BuildReplayReport(nil)
	if r.FrameCount != 0 {
		t.Errorf("expected FrameCount=0 for empty input")
	}
	if len(r.Frames) != 0 {
		t.Errorf("expected empty Frames slice")
	}
}
