package snapshot

import (
	"path/filepath"
	"testing"
)

func TestAnnotationStore_SetAndGet(t *testing.T) {
	store := NewAnnotationStore(filepath.Join(t.TempDir(), "ann.json"))
	as, _ := NewAnnotationSet([]string{"env=prod"})
	store.Set("tcp", 8080, as)

	got := store.Get("tcp", 8080)
	v, ok := got.Get("env")
	if !ok || v != "prod" {
		t.Fatalf("expected env=prod, got %q ok=%v", v, ok)
	}
}

func TestAnnotationStore_Get_Missing(t *testing.T) {
	store := NewAnnotationStore(filepath.Join(t.TempDir(), "ann.json"))
	got := store.Get("udp", 53)
	if got.Len() != 0 {
		t.Fatal("expected empty AnnotationSet for missing key")
	}
}

func TestAnnotationStore_Set_NilClearsKey(t *testing.T) {
	store := NewAnnotationStore(filepath.Join(t.TempDir(), "ann.json"))
	as, _ := NewAnnotationSet([]string{"x=y"})
	store.Set("tcp", 443, as)
	store.Set("tcp", 443, nil)
	if len(store.Keys()) != 0 {
		t.Fatal("expected key to be removed after Set with nil")
	}
}

func TestAnnotationStore_SaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "ann.json")
	store := NewAnnotationStore(path)
	as, _ := NewAnnotationSet([]string{"owner=bob", "tier=backend"})
	store.Set("tcp", 9090, as)

	if err := store.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	store2 := NewAnnotationStore(path)
	if err := store2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	got := store2.Get("tcp", 9090)
	v, ok := got.Get("owner")
	if !ok || v != "bob" {
		t.Fatalf("expected owner=bob after reload, got %q ok=%v", v, ok)
	}
}

func TestAnnotationStore_Load_MissingFile(t *testing.T) {
	store := NewAnnotationStore(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err := store.Load(); err != nil {
		t.Fatalf("Load of missing file should not error, got: %v", err)
	}
}

func TestAnnotationStore_Keys(t *testing.T) {
	store := NewAnnotationStore(filepath.Join(t.TempDir(), "ann.json"))
	a1, _ := NewAnnotationSet([]string{"a=1"})
	a2, _ := NewAnnotationSet([]string{"b=2"})
	store.Set("tcp", 80, a1)
	store.Set("udp", 53, a2)
	keys := store.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}
