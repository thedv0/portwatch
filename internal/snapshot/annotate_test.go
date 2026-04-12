package snapshot

import (
	"testing"
)

func TestNewAnnotationSet_Valid(t *testing.T) {
	as, err := NewAnnotationSet([]string{"env=prod", "owner=alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if as.Len() != 2 {
		t.Fatalf("expected 2 annotations, got %d", as.Len())
	}
}

func TestNewAnnotationSet_InvalidFormat(t *testing.T) {
	_, err := NewAnnotationSet([]string{"badvalue"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewAnnotationSet_EmptyKey(t *testing.T) {
	_, err := NewAnnotationSet([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestAnnotationSet_GetAndAdd(t *testing.T) {
	as, _ := NewAnnotationSet(nil)
	as.Add("team", "platform")
	v, ok := as.Get("team")
	if !ok || v != "platform" {
		t.Fatalf("expected 'platform', got %q ok=%v", v, ok)
	}
}

func TestAnnotationSet_Add_Overwrites(t *testing.T) {
	as, _ := NewAnnotationSet([]string{"env=staging"})
	as.Add("env", "prod")
	v, _ := as.Get("env")
	if v != "prod" {
		t.Fatalf("expected overwritten value 'prod', got %q", v)
	}
	if as.Len() != 1 {
		t.Fatalf("expected len 1 after overwrite, got %d", as.Len())
	}
}

func TestAnnotationSet_Get_Missing(t *testing.T) {
	as, _ := NewAnnotationSet(nil)
	_, ok := as.Get("nonexistent")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestAnnotationSet_Remove(t *testing.T) {
	as, _ := NewAnnotationSet([]string{"a=1", "b=2"})
	removed := as.Remove("a")
	if !removed {
		t.Fatal("expected Remove to return true")
	}
	if as.Len() != 1 {
		t.Fatalf("expected len 1 after remove, got %d", as.Len())
	}
	_, ok := as.Get("a")
	if ok {
		t.Fatal("expected key 'a' to be gone")
	}
}

func TestAnnotationSet_Remove_Missing(t *testing.T) {
	as, _ := NewAnnotationSet(nil)
	if as.Remove("x") {
		t.Fatal("expected Remove to return false for missing key")
	}
}

func TestAnnotationSet_All_ReturnsCopy(t *testing.T) {
	as, _ := NewAnnotationSet([]string{"k=v"})
	all := as.All()
	all[0].Value = "mutated"
	v, _ := as.Get("k")
	if v != "v" {
		t.Fatal("All() should return a copy, not a reference")
	}
}
