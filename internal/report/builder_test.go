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
