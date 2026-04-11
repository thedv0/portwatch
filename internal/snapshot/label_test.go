package snapshot

import (
	"testing"
)

func TestNewLabelSet_Valid(t *testing.T) {
	ls, err := NewLabelSet([]string{"env=prod", "team=infra"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ls.Len() != 2 {
		t.Fatalf("expected 2 labels, got %d", ls.Len())
	}
}

func TestNewLabelSet_InvalidFormat(t *testing.T) {
	_, err := NewLabelSet([]string{"badlabel"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewLabelSet_EmptyKey(t *testing.T) {
	_, err := NewLabelSet([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestLabelSet_SetAndGet(t *testing.T) {
	ls, _ := NewLabelSet(nil)
	if err := ls.Set("region", "us-east"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	v, ok := ls.Get("region")
	if !ok || v != "us-east" {
		t.Fatalf("expected us-east, got %q ok=%v", v, ok)
	}
}

func TestLabelSet_Set_EmptyKey(t *testing.T) {
	ls, _ := NewLabelSet(nil)
	if err := ls.Set("", "value"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestLabelSet_Delete(t *testing.T) {
	ls, _ := NewLabelSet([]string{"a=1"})
	ls.Delete("a")
	_, ok := ls.Get("a")
	if ok {
		t.Fatal("expected label to be deleted")
	}
}

func TestLabelSet_All_Sorted(t *testing.T) {
	ls, _ := NewLabelSet([]string{"z=last", "a=first", "m=mid"})
	all := ls.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(all))
	}
	if all[0].Key != "a" || all[1].Key != "m" || all[2].Key != "z" {
		t.Fatalf("labels not sorted: %v", all)
	}
}

func TestLabelSet_Merge(t *testing.T) {
	ls1, _ := NewLabelSet([]string{"a=1", "b=2"})
	ls2, _ := NewLabelSet([]string{"b=override", "c=3"})
	ls1.Merge(ls2)
	if ls1.Len() != 3 {
		t.Fatalf("expected 3 labels after merge, got %d", ls1.Len())
	}
	v, _ := ls1.Get("b")
	if v != "override" {
		t.Fatalf("expected b=override, got %q", v)
	}
}

func TestLabelSet_Get_Missing(t *testing.T) {
	ls, _ := NewLabelSet(nil)
	_, ok := ls.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}
