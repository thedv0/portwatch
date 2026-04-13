package snapshot

import (
	"testing"
	"time"
)

func makeInterpSnap(ts time.Time, ports []Port) Snapshot {
	return Snapshot{Timestamp: ts, Ports: ports}
}

func TestInterpolate_EmptyInput(t *testing.T) {
	res, err := Interpolate(nil, time.Minute, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Snapshots) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(res.Snapshots))
	}
}

func TestInterpolate_InvalidStep(t *testing.T) {
	snaps := []Snapshot{makeInterpSnap(time.Now(), nil)}
	_, err := Interpolate(snaps, 0, DefaultInterpolateOptions())
	if err == nil {
		t.Fatal("expected error for zero step")
	}
}

func TestInterpolate_FillsGap_ForwardMethod(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	ports := []Port{{Port: 80, Protocol: "tcp"}}
	snaps := []Snapshot{
		makeInterpSnap(now, ports),
		makeInterpSnap(now.Add(3*time.Minute), nil),
	}
	opts := DefaultInterpolateOptions()
	opts.Method = "forward"
	res, err := Interpolate(snaps, time.Minute, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Filled != 2 {
		t.Errorf("expected 2 filled snapshots, got %d", res.Filled)
	}
	if len(res.Snapshots) != 5 {
		t.Errorf("expected 5 total snapshots, got %d", len(res.Snapshots))
	}
	for _, s := range res.Snapshots[1:3] {
		if len(s.Ports) != 1 || s.Ports[0].Port != 80 {
			t.Errorf("forward fill: expected port 80, got %+v", s.Ports)
		}
	}
}

func TestInterpolate_ZeroMethod_EmptyPorts(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	ports := []Port{{Port: 443, Protocol: "tcp"}}
	snaps := []Snapshot{
		makeInterpSnap(now, ports),
		makeInterpSnap(now.Add(2*time.Minute), nil),
	}
	opts := DefaultInterpolateOptions()
	opts.Method = "zero"
	res, err := Interpolate(snaps, time.Minute, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Filled != 1 {
		t.Errorf("expected 1 filled snapshot, got %d", res.Filled)
	}
	for _, s := range res.Snapshots[1:2] {
		if len(s.Ports) != 0 {
			t.Errorf("zero fill: expected empty ports, got %+v", s.Ports)
		}
	}
}

func TestInterpolate_GapExceedsMaxGap_NoFill(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	snaps := []Snapshot{
		makeInterpSnap(now, nil),
		makeInterpSnap(now.Add(30*time.Minute), nil),
	}
	opts := DefaultInterpolateOptions()
	opts.MaxGap = 5 * time.Minute
	res, err := Interpolate(snaps, time.Minute, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Filled != 0 {
		t.Errorf("expected 0 filled snapshots for gap > MaxGap, got %d", res.Filled)
	}
	if len(res.Snapshots) != 2 {
		t.Errorf("expected 2 snapshots unchanged, got %d", len(res.Snapshots))
	}
}

func TestInterpolate_DefaultOptions_Values(t *testing.T) {
	opts := DefaultInterpolateOptions()
	if opts.Method != "forward" {
		t.Errorf("expected method 'forward', got %q", opts.Method)
	}
	if opts.MaxGap != 5*time.Minute {
		t.Errorf("expected MaxGap 5m, got %v", opts.MaxGap)
	}
}
