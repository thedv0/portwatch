package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultRetentionPolicy_Values(t *testing.T) {
	p := DefaultRetentionPolicy()
	if p.MaxAge != 7*24*time.Hour {
		t.Errorf("expected MaxAge=7d, got %v", p.MaxAge)
	}
	if p.MaxCount != 100 {
		t.Errorf("expected MaxCount=100, got %d", p.MaxCount)
	}
	if p.MaxHistoryEntries != 500 {
		t.Errorf("expected MaxHistoryEntries=500, got %d", p.MaxHistoryEntries)
	}
}

func TestValidate_Valid(t *testing.T) {
	p := DefaultRetentionPolicy()
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_NegativeMaxAge(t *testing.T) {
	p := DefaultRetentionPolicy()
	p.MaxAge = -1
	if err := p.Validate(); err == nil {
		t.Error("expected error for negative MaxAge")
	}
}

func TestValidate_NegativeMaxCount(t *testing.T) {
	p := DefaultRetentionPolicy()
	p.MaxCount = -5
	if err := p.Validate(); err == nil {
		t.Error("expected error for negative MaxCount")
	}
}

func TestValidate_NegativeMaxHistoryEntries(t *testing.T) {
	p := DefaultRetentionPolicy()
	p.MaxHistoryEntries = -1
	if err := p.Validate(); err == nil {
		t.Error("expected error for negative MaxHistoryEntries")
	}
}

func TestApply_RemovesOldFiles(t *testing.T) {
	dir := t.TempDir()

	// Write an old file
	old := filepath.Join(dir, "snap_old.json")
	if err := os.WriteFile(old, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().Add(-10 * 24 * time.Hour)
	if err := os.Chtimes(old, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	p := RetentionPolicy{
		MaxAge:            7 * 24 * time.Hour,
		MaxCount:          100,
		MaxHistoryEntries: 500,
	}
	if err := p.Apply(dir); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("expected old snapshot file to be removed")
	}
}
