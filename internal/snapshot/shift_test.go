package snapshot

import (
	"testing"
	"time"
)

func makeShiftSnap(ts time.Time) Snapshot {
	return Snapshot{
		Timestamp: ts,
		Ports:     []Port{{Port: 80, Protocol: "tcp", PID: 1}},
	}
}

func TestShift_PositiveOffset(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	snaps := []Snapshot{makeShiftSnap(base)}

	opts := DefaultShiftOptions()
	opts.Offset = 2 * time.Hour
	opts.Clock = func() time.Time { return base.Add(24 * time.Hour) }

	out, err := Shift(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(out))
	}
	want := base.Add(2 * time.Hour)
	if !out[0].Timestamp.Equal(want) {
		t.Errorf("expected %v, got %v", want, out[0].Timestamp)
	}
}

func TestShift_NegativeOffset(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	snaps := []Snapshot{makeShiftSnap(base)}

	opts := DefaultShiftOptions()
	opts.Offset = -30 * time.Minute
	opts.Clock = time.Now

	out, err := Shift(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := base.Add(-30 * time.Minute)
	if !out[0].Timestamp.Equal(want) {
		t.Errorf("expected %v, got %v", want, out[0].Timestamp)
	}
}

func TestShift_ClampToNow(t *testing.T) {
	now := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	future := now.Add(5 * time.Hour)
	snaps := []Snapshot{makeShiftSnap(future)}

	opts := DefaultShiftOptions()
	opts.Offset = 2 * time.Hour
	opts.ClampToNow = true
	opts.Clock = func() time.Time { return now }

	out, err := Shift(snaps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out[0].Timestamp.Equal(now) {
		t.Errorf("expected clamped timestamp %v, got %v", now, out[0].Timestamp)
	}
}

func TestShift_NilClock_ReturnsError(t *testing.T) {
	opts := DefaultShiftOptions()
	opts.Clock = nil

	_, err := Shift([]Snapshot{makeShiftSnap(time.Now())}, opts)
	if err == nil {
		t.Fatal("expected error for nil Clock, got nil")
	}
}

func TestShift_EmptyInput_ReturnsEmpty(t *testing.T) {
	opts := DefaultShiftOptions()
	opts.Clock = time.Now

	out, err := Shift([]Snapshot{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d snapshots", len(out))
	}
}

func TestShift_PreservesPortData(t *testing.T) {
	base := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	snap := makeShiftSnap(base)
	snap.Ports = append(snap.Ports, Port{Port: 443, Protocol: "tcp", PID: 42})

	opts := DefaultShiftOptions()
	opts.Offset = time.Hour
	opts.Clock = time.Now

	out, err := Shift([]Snapshot{snap}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out[0].Ports) != 2 {
		t.Errorf("expected 2 ports preserved, got %d", len(out[0].Ports))
	}
}
