package snapshot

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// CleanerConfig holds configuration for snapshot cleanup.
type CleanerConfig struct {
	// MaxAge is the maximum age of snapshot files to retain.
	MaxAge time.Duration
	// MaxFiles is the maximum number of snapshot files to retain.
	MaxFiles int
	// Dir is the directory containing snapshot files.
	Dir string
}

// Cleaner removes old snapshot files based on configured retention policy.
type Cleaner struct {
	cfg CleanerConfig
}

// NewCleaner creates a new Cleaner with the given configuration.
func NewCleaner(cfg CleanerConfig) *Cleaner {
	return &Cleaner{cfg: cfg}
}

// Clean removes snapshot files that exceed the retention policy.
// Files are removed if they are older than MaxAge or if there are more
// than MaxFiles files, removing the oldest first.
func (c *Cleaner) Clean() error {
	entries, err := os.ReadDir(c.cfg.Dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var snaps []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			snaps = append(snaps, e)
		}
	}

	sort.Slice(snaps, func(i, j int) bool {
		ii, _ := snaps[i].Info()
		jj, _ := snaps[j].Info()
		if ii == nil || jj == nil {
			return false
		}
		return ii.ModTime().Before(jj.ModTime())
	})

	now := time.Now()
	for i, e := range snaps {
		info, err := e.Info()
		if err != nil {
			continue
		}
		exceedsAge := c.cfg.MaxAge > 0 && now.Sub(info.ModTime()) > c.cfg.MaxAge
		exceedsCount := c.cfg.MaxFiles > 0 && i < len(snaps)-c.cfg.MaxFiles
		if exceedsAge || exceedsCount {
			path := filepath.Join(c.cfg.Dir, e.Name())
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

// Count returns the number of snapshot files currently in the configured directory.
func (c *Cleaner) Count() (int, error) {
	entries, err := os.ReadDir(c.cfg.Dir)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			count++
		}
	}
	return count, nil
}
