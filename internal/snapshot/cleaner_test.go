package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeSnapFile(t *testing.T, dir, name string, modTime time.Time) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write snap file: %v", err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
}

func TestClean_MissingDirIsNoop(t *testing.T) {
	c := NewCleaner(CleanerConfig{Dir: "/nonexistent/path/portwatch"})
	if err := c.Clean(); err != nil {
		t.Fatalf("expected no error for missing dir, got: %v", err)
	}
}

func TestClean_RemovesByAge(t *testing.T) {
	dir := t.TempDir()
	old := time.Now().Add(-48 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)

	writeSnapFile(t, dir, "old.json", old)
	writeSnapFile(t, dir, "recent.json", recent)

	c := NewCleaner(CleanerConfig{Dir: dir, MaxAge: 24 * time.Hour})
	if err := c.Clean(); err != nil {
		t.Fatalf("clean: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "old.json")); !os.IsNotExist(err) {
		t.Error("expected old.json to be removed")
	}
	if _, err := os.Stat(filepath.Join(dir, "recent.json")); err != nil {
		t.Errorf("expected recent.json to exist: %v", err)
	}
}

func TestClean_RemovesByCount(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()

	for i, name := range []string{"a.json", "b.json", "c.json", "d.json"} {
		writeSnapFile(t, dir, name, now.Add(time.Duration(i)*time.Minute))
	}

	c := NewCleaner(CleanerConfig{Dir: dir, MaxFiles: 2})
	if err := c.Clean(); err != nil {
		t.Fatalf("clean: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 2 {
		t.Errorf("expected 2 files remaining, got %d", len(entries))
	}
}

func TestClean_NoFilesIsNoop(t *testing.T) {
	dir := t.TempDir()
	c := NewCleaner(CleanerConfig{Dir: dir, MaxAge: time.Hour, MaxFiles: 10})
	if err := c.Clean(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClean_IgnoresNonJSONFiles(t *testing.T) {
	dir := t.TempDir()
	old := time.Now().Add(-72 * time.Hour)
	writeSnapFile(t, dir, "snap.json", old)

	other := filepath.Join(dir, "notes.txt")
	_ = os.WriteFile(other, []byte("hello"), 0o644)
	_ = os.Chtimes(other, old, old)

	c := NewCleaner(CleanerConfig{Dir: dir, MaxFiles: 0, MaxAge: time.Hour})
	if err := c.Clean(); err != nil {
		t.Fatalf("clean: %v", err)
	}

	if _, err := os.Stat(other); err != nil {
		t.Error("expected notes.txt to be untouched")
	}
}
