package snapshot

import (
	"testing"
	"time"
)

func makeReplaySnap(t time.Time, ports []PortState) Snapshot {
	return Snapshot{Timestamp: t, Ports: ports}
}

func TestReplay_EmptyInput(t *testing.T) {
	frames, err := Replay(nil, DefaultReplayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) != 0 {
		t.Fatalf("expected 0 frames, got %d", len(frames))
	}
}

func TestReplay_FrameCount(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeReplaySnap(now, []PortState{{Port: 80, Protocol: "tcp"}}),
		makeReplaySnap(now.Add(time.Minute), []PortState{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}}),
		makeReplaySnap(now.Add(2*time.Minute), []PortState{{Port: 443, Protocol: "tcp"}}),
	}
	frames, err := Replay(snaps, DefaultReplayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) != 3 {
		t.Fatalf("expected 3 frames, got %d", len(frames))
	}
}

func TestReplay_DiffPopulated(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeReplaySnap(now, []PortState{{Port: 80, Protocol: "tcp"}}),
		makeReplaySnap(now.Add(time.Minute), []PortState{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}}),
	}
	frames, _ := Replay(snaps, DefaultReplayOptions())
	if len(frames[0].Diff.Added) != 1 {
		t.Errorf("first frame should have 1 added port")
	}
	if len(frames[1].Diff.Added) != 1 {
		t.Errorf("second frame should have 1 added port (443)")
	}
}

func TestReplay_MaxFrames(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeReplaySnap(now, []PortState{}),
		makeReplaySnap(now.Add(time.Minute), []PortState{}),
		makeReplaySnap(now.Add(2*time.Minute), []PortState{}),
	}
	opts := DefaultReplayOptions()
	opts.MaxFrames = 2
	frames, _ := Replay(snaps, opts)
	if len(frames) != 2 {
		t.Fatalf("expected 2 frames, got %d", len(frames))
	}
}

func TestReplay_TimeFilter(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeReplaySnap(now, []PortState{}),
		makeReplaySnap(now.Add(time.Minute), []PortState{}),
		makeReplaySnap(now.Add(2*time.Minute), []PortState{}),
	}
	opts := DefaultReplayOptions()
	opts.StartAt = now.Add(30 * time.Second)
	frames, _ := Replay(snaps, opts)
	if len(frames) != 2 {
		t.Fatalf("expected 2 frames after time filter, got %d", len(frames))
	}
}

func TestReplay_Reverse(t *testing.T) {
	now := time.Now()
	snaps := []Snapshot{
		makeReplaySnap(now, []PortState{}),
		makeReplaySnap(now.Add(time.Minute), []PortState{}),
	}
	opts := DefaultReplayOptions()
	opts.Reverse = true
	frames, _ := Replay(snaps, opts)
	if !frames[0].Timestamp.Equal(now.Add(time.Minute)) {
		t.Errorf("expected newest frame first in reverse mode")
	}
}
