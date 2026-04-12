package snapshot

import (
	"path/filepath"
	"testing"
)

func TestPinStore_PinAndLoad(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(filepath.Join(dir, "pins.json"))

	err := s.Pin(PinnedPort{Port: 8080, Protocol: "tcp", Comment: "web server"})
	if err != nil {
		t.Fatalf("Pin: %v", err)
	}

	pins, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(pins))
	}
	p := pins["tcp:8080"]
	if p.Port != 8080 || p.Protocol != "tcp" || p.Comment != "web server" {
		t.Errorf("unexpected pin: %+v", p)
	}
	if p.PinnedAt.IsZero() {
		t.Error("PinnedAt should be set")
	}
}

func TestPinStore_IsPinned(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(filepath.Join(dir, "pins.json"))

	_ = s.Pin(PinnedPort{Port: 443, Protocol: "tcp"})

	ok, err := s.IsPinned(443, "tcp")
	if err != nil || !ok {
		t.Errorf("expected 443/tcp to be pinned, got ok=%v err=%v", ok, err)
	}

	ok, err = s.IsPinned(80, "tcp")
	if err != nil || ok {
		t.Errorf("expected 80/tcp not to be pinned, got ok=%v err=%v", ok, err)
	}
}

func TestPinStore_Unpin(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(filepath.Join(dir, "pins.json"))

	_ = s.Pin(PinnedPort{Port: 22, Protocol: "tcp"})
	_ = s.Pin(PinnedPort{Port: 8080, Protocol: "tcp"})

	if err := s.Unpin(22, "tcp"); err != nil {
		t.Fatalf("Unpin: %v", err)
	}

	ok, _ := s.IsPinned(22, "tcp")
	if ok {
		t.Error("expected 22/tcp to be unpinned")
	}
	ok, _ = s.IsPinned(8080, "tcp")
	if !ok {
		t.Error("expected 8080/tcp to remain pinned")
	}
}

func TestPinStore_Load_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(filepath.Join(dir, "nonexistent", "pins.json"))

	pins, err := s.Load()
	if err != nil {
		t.Fatalf("Load on missing file: %v", err)
	}
	if len(pins) != 0 {
		t.Errorf("expected empty map, got %d entries", len(pins))
	}
}

func TestPinStore_Pin_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(filepath.Join(dir, "nested", "deep", "pins.json"))

	if err := s.Pin(PinnedPort{Port: 9090, Protocol: "udp"}); err != nil {
		t.Fatalf("Pin with nested path: %v", err)
	}
	ok, _ := s.IsPinned(9090, "udp")
	if !ok {
		t.Error("expected 9090/udp to be pinned after nested save")
	}
}
