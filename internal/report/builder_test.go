package report_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/snapshot"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestBuild_SetsTimestamp(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	b := report.NewBuilder().WithClock(fixedClock(ts))
	r := b.Build(snapshot.DiffResult{}, 3)
	if !r.Timestamp.Equal(ts) {
		t.Errorf("expected timestamp %v, got %v", ts, r.Timestamp)
	}
}

func TestBuild_TotalOpen(t *testing.T) {
	b := report.NewBuilder()
	r := b.Build(snapshot.DiffResult{}, 42)
	if r.Total != 42 {
		t.Errorf("expected Total=42, got %d", r.Total)
	}
}

func TestBuild_DiffPopulated(t *testing.T) {
	added := []snapshot.Port{{Port: 443, Protocol: "tcp", PID: 99}}
	removed := []snapshot.Port{{Port: 80, Protocol: "tcp", PID: 10}}
	diff := snapshot.DiffResult{Added: added, Removed: removed}

	b := report.NewBuilder()
	r := b.Build(diff, 10)

	if len(r.Added) != 1 || r.Added[0].Port != 443 {
		t.Errorf("unexpected Added: %+v", r.Added)
	}
	if len(r.Removed) != 1 || r.Removed[0].Port != 80 {
		t.Errorf("unexpected Removed: %+v", r.Removed)
	}
}

func TestBuild_EmptyDiff(t *testing.T) {
	b := report.NewBuilder()
	r := b.Build(snapshot.DiffResult{}, 0)
	if len(r.Added) != 0 || len(r.Removed) != 0 {
		t.Errorf("expected empty slices for no-change diff")
	}
}

func TestNewBuilder_DefaultClockIsRecent(t *testing.T) {
	before := time.Now()
	b := report.NewBuilder()
	r := b.Build(snapshot.DiffResult{}, 0)
	after := time.Now()
	if r.Timestamp.Before(before) || r.Timestamp.After(after) {
		t.Errorf("timestamp %v not between %v and %v", r.Timestamp, before, after)
	}
}

func TestBuild_MultipleCalls_IndependentTimestamps(t *testing.T) {
	ts1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	ts2 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	b1 := report.NewBuilder().WithClock(fixedClock(ts1))
	b2 := report.NewBuilder().WithClock(fixedClock(ts2))

	r1 := b1.Build(snapshot.DiffResult{}, 5)
	r2 := b2.Build(snapshot.DiffResult{}, 10)

	if !r1.Timestamp.Equal(ts1) {
		t.Errorf("b1: expected timestamp %v, got %v", ts1, r1.Timestamp)
	}
	if !r2.Timestamp.Equal(ts2) {
		t.Errorf("b2: expected timestamp %v, got %v", ts2, r2.Timestamp)
	}
}
