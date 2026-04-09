package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Protocol: "tcp", Process: "nginx"},
		{Port: 443, Protocol: "tcp", Process: "nginx"},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	store := snapshot.NewStore(path)

	ports := makePorts()
	if err := store.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	state, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(state.Ports) != len(ports) {
		t.Errorf("got %d ports, want %d", len(state.Ports), len(ports))
	}
	if state.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	store := snapshot.NewStore("/tmp/portwatch-nonexistent-snap-xyz.json")
	state, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !state.Timestamp.IsZero() || len(state.Ports) != 0 {
		t.Error("expected empty state for missing file")
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "snap.json")
	store := snapshot.NewStore(path)

	if err := store.Save(makePorts()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file at %s: %v", path, err)
	}
}

func TestSave_TimestampIsRecent(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(filepath.Join(dir, "snap.json"))
	before := time.Now().UTC().Add(-time.Second)

	if err := store.Save(makePorts()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	state, _ := store.Load()
	if state.Timestamp.Before(before) {
		t.Errorf("timestamp %v is before test start %v", state.Timestamp, before)
	}
}
