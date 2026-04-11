package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultRotationPolicy_Values(t *testing.T) {
	p := DefaultRotationPolicy()
	if p.MaxSizeBytes != 10*1024*1024 {
		t.Errorf("expected 10MB, got %d", p.MaxSizeBytes)
	}
	if p.MaxAgeDays != 7 {
		t.Errorf("expected 7 days, got %d", p.MaxAgeDays)
	}
	if p.MaxBackups != 5 {
		t.Errorf("expected 5 backups, got %d", p.MaxBackups)
	}
}

func TestValidate_RotationPolicy_Valid(t *testing.T) {
	if err := DefaultRotationPolicy().Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_RotationPolicy_Negative(t *testing.T) {
	cases := []RotationPolicy{
		{MaxSizeBytes: -1},
		{MaxAgeDays: -1},
		{MaxBackups: -1},
	}
	for _, p := range cases {
		if err := p.Validate(); err == nil {
			t.Errorf("expected error for policy %+v", p)
		}
	}
}

func TestNeedsRotation_SizeExceeded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	if err := os.WriteFile(path, []byte(strings.Repeat("x", 100)), 0o644); err != nil {
		t.Fatal(err)
	}
	p := RotationPolicy{MaxSizeBytes: 50}
	need, err := NeedsRotation(path, p)
	if err != nil {
		t.Fatal(err)
	}
	if !need {
		t.Error("expected rotation needed")
	}
}

func TestNeedsRotation_SizeNotExceeded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	if err := os.WriteFile(path, []byte("small"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := RotationPolicy{MaxSizeBytes: 1024}
	need, err := NeedsRotation(path, p)
	if err != nil {
		t.Fatal(err)
	}
	if need {
		t.Error("expected no rotation")
	}
}

func TestNeedsRotation_MissingFile(t *testing.T) {
	need, err := NeedsRotation("/nonexistent/audit.log", DefaultRotationPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if need {
		t.Error("expected no rotation for missing file")
	}
}

func TestRotate_RenamesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	if err := os.WriteFile(path, []byte("log data"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := RotationPolicy{MaxBackups: 10}
	if err := Rotate(path, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("original file should have been renamed")
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 backup file, got %d", len(entries))
	}
	if !strings.HasPrefix(entries[0].Name(), "audit.") {
		t.Errorf("unexpected backup name: %s", entries[0].Name())
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	// Create existing backups
	for i := 0; i < 5; i++ {
		ts := time.Now().UTC().Add(time.Duration(-i) * time.Minute).Format("20060102T150405Z")
		name := filepath.Join(dir, "audit."+ts+".log")
		_ = os.WriteFile(name, []byte("old"), 0o644)
	}

	if err := os.WriteFile(path, []byte("current"), 0o644); err != nil {
		t.Fatal(err)
	}

	p := RotationPolicy{MaxBackups: 3}
	if err := Rotate(path, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 3 {
		t.Errorf("expected 3 backups after prune, got %d", len(entries))
	}
}
