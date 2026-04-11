package snapshot

import (
	"testing"
)

func TestNewTagSet_Valid(t *testing.T) {
	ts, err := NewTagSet([]string{"env=prod", "region=us-east"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := ts.Get("env"); !ok || v != "prod" {
		t.Errorf("expected env=prod, got %q ok=%v", v, ok)
	}
	if v, ok := ts.Get("region"); !ok || v != "us-east" {
		t.Errorf("expected region=us-east, got %q ok=%v", v, ok)
	}
}

func TestNewTagSet_NoValue(t *testing.T) {
	ts, err := NewTagSet([]string{"debug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ts.Has("debug") {
		t.Error("expected 'debug' key to be present")
	}
	if v, _ := ts.Get("debug"); v != "" {
		t.Errorf("expected empty value, got %q", v)
	}
}

func TestNewTagSet_EmptyKey(t *testing.T) {
	_, err := NewTagSet([]string{"=value"})
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestTagSet_Has_Missing(t *testing.T) {
	ts, _ := NewTagSet([]string{"a=1"})
	if ts.Has("b") {
		t.Error("expected Has('b') to be false")
	}
}

func TestTagSet_Merge(t *testing.T) {
	base, _ := NewTagSet([]string{"a=1", "b=2"})
	over, _ := NewTagSet([]string{"b=99", "c=3"})
	merged := base.Merge(over)

	if v, _ := merged.Get("a"); v != "1" {
		t.Errorf("expected a=1, got %q", v)
	}
	if v, _ := merged.Get("b"); v != "99" {
		t.Errorf("expected b=99 (override), got %q", v)
	}
	if v, _ := merged.Get("c"); v != "3" {
		t.Errorf("expected c=3, got %q", v)
	}
}

func TestTagSet_Slice_Sorted(t *testing.T) {
	ts, _ := NewTagSet([]string{"z=last", "a=first", "m=mid"})
	slice := ts.Slice()
	if len(slice) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(slice))
	}
	if slice[0] != "a=first" || slice[1] != "m=mid" || slice[2] != "z=last" {
		t.Errorf("unexpected order: %v", slice)
	}
}

func TestNewTagSet_Empty(t *testing.T) {
	ts, err := NewTagSet(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ts) != 0 {
		t.Errorf("expected empty TagSet, got %d entries", len(ts))
	}
}
