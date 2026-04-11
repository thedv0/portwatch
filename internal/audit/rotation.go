package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RotationPolicy defines when and how audit log files are rotated.
type RotationPolicy struct {
	MaxSizeBytes int64  // rotate when file exceeds this size (0 = disabled)
	MaxAgeDays   int    // rotate files older than this many days (0 = disabled)
	MaxBackups   int    // maximum number of rotated files to keep (0 = keep all)
}

// DefaultRotationPolicy returns a sensible default rotation policy.
func DefaultRotationPolicy() RotationPolicy {
	return RotationPolicy{
		MaxSizeBytes: 10 * 1024 * 1024, // 10 MB
		MaxAgeDays:   7,
		MaxBackups:   5,
	}
}

// Validate checks that the rotation policy fields are non-negative.
func (p RotationPolicy) Validate() error {
	if p.MaxSizeBytes < 0 {
		return fmt.Errorf("rotation: MaxSizeBytes must be >= 0")
	}
	if p.MaxAgeDays < 0 {
		return fmt.Errorf("rotation: MaxAgeDays must be >= 0")
	}
	if p.MaxBackups < 0 {
		return fmt.Errorf("rotation: MaxBackups must be >= 0")
	}
	return nil
}

// Rotate renames the current log file with a timestamp suffix and removes
// old backups according to the policy.
func Rotate(path string, policy RotationPolicy) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	dest := fmt.Sprintf("%s.%s%s", base, timestamp, ext)

	if err := os.Rename(path, dest); err != nil {
		return fmt.Errorf("rotate: rename %s -> %s: %w", path, dest, err)
	}

	if policy.MaxBackups > 0 {
		if err := pruneBackups(path, policy.MaxBackups); err != nil {
			return err
		}
	}
	return nil
}

// NeedsRotation reports whether the file at path should be rotated given the policy.
func NeedsRotation(path string, policy RotationPolicy) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if policy.MaxSizeBytes > 0 && info.Size() >= policy.MaxSizeBytes {
		return true, nil
	}
	if policy.MaxAgeDays > 0 {
		cutoff := time.Now().UTC().AddDate(0, 0, -policy.MaxAgeDays)
		if info.ModTime().Before(cutoff) {
			return true, nil
		}
	}
	return false, nil
}

func pruneBackups(originalPath string, maxBackups int) error {
	dir := filepath.Dir(originalPath)
	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(filepath.Base(originalPath), ext)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("rotate: read dir %s: %w", dir, err)
	}

	var backups []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, base+".") && strings.HasSuffix(name, ext) && name != filepath.Base(originalPath) {
			backups = append(backups, filepath.Join(dir, name))
		}
	}

	sort.Strings(backups)
	for len(backups) > maxBackups {
		if err := os.Remove(backups[0]); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rotate: remove old backup %s: %w", backups[0], err)
		}
		backups = backups[1:]
	}
	return nil
}
