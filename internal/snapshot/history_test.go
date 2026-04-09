package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendHistory_CreatesAndLoads(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	entry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Added:     []string{"tcp:8080"},
	}

	if err := AppendHistory(path, entry); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 1 || len(entries[0].Added) != 1 {
		t.Fatalf("unexpected entries: %v", entries)
	}
}

func TestLoadHistory_MissingFile(t *testing.T) {
	entries, err := LoadHistory("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Fatalf("expected nil entries, got %v", entries)
	}
}

func TestAppendHistory_TrimsToMax(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	base := time.Now().UTC()
	for i := 0; i < maxHistoryEntries+10; i++ {
		e := HistoryEntry{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Added:     []string{"tcp:9000"},
		}
		if err := AppendHistory(path, e); err != nil {
			t.Fatalf("AppendHistory iteration %d: %v", i, err)
		}
	}

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != maxHistoryEntries {
		t.Fatalf("expected %d entries after trim, got %d", maxHistoryEntries, len(entries))
	}
}

func TestAppendHistory_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", "history.json")

	if err := AppendHistory(path, HistoryEntry{Timestamp: time.Now()}); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}
