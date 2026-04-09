package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const maxHistoryEntries = 48

// HistoryEntry records a diff event with a timestamp.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Added     []string  `json:"added,omitempty"`
	Removed   []string  `json:"removed,omitempty"`
}

// AppendHistory loads the history file, appends the entry, trims to max, and saves.
func AppendHistory(path string, entry HistoryEntry) error {
	entries, _ := LoadHistory(path)
	entries = append(entries, entry)

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	if len(entries) > maxHistoryEntries {
		entries = entries[len(entries)-maxHistoryEntries:]
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("history mkdir: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("history create: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

// LoadHistory reads history entries from disk; returns empty slice if missing.
func LoadHistory(path string) ([]HistoryEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []HistoryEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("history decode: %w", err)
	}
	return entries, nil
}
