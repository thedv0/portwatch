package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeBaselinePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Protocol: "tcp", PID: 100, Process: "nginx"},
		{Port: 443, Protocol: "tcp", PID: 101, Process: "nginx"},
	}
}

func TestBaselineStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewBaselineStore(dir)
	ports := makeBaselinePorts()

	if err := store.Save("default", ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := store.Load("default")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b.Ports) != len(ports) {
		t.Errorf("ports len: got %d, want %d", len(b.Ports), len(ports))
	}
	if b.Label != "default" {
		t.Errorf("label: got %q, want %q", b.Label, "default")
	}
}

func TestBaselineStore_Load_Missing(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewBaselineStore(dir)

	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestBaselineStore_Save_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "baselines")
	store := snapshot.NewBaselineStore(dir)

	if err := store.Save("init", makeBaselinePorts()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("dir not created: %v", err)
	}
}

func TestBaselineStore_List(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewBaselineStore(dir)

	for _, label := range []string{"prod", "staging", "dev"} {
		if err := store.Save(label, makeBaselinePorts()); err != nil {
			t.Fatalf("Save %s: %v", label, err)
		}
	}

	labels, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 3 {
		t.Errorf("list len: got %d, want 3", len(labels))
	}
}

func TestBaselineStore_List_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewBaselineStore(dir)

	labels, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty list, got %v", labels)
	}
}

func TestBaselineStore_TimestampIsRecent(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewBaselineStore(dir)
	before := time.Now().UTC().Add(-time.Second)

	if err := store.Save("ts", makeBaselinePorts()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := store.Load("ts")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if b.CreatedAt.Before(before) {
		t.Errorf("timestamp too old: %v", b.CreatedAt)
	}
}
