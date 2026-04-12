package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PinnedPort represents a port that has been explicitly pinned (marked as known-good).
type PinnedPort struct {
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Process   string    `json:"process,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	PinnedAt  time.Time `json:"pinned_at"`
}

// PinStore persists pinned ports to disk.
type PinStore struct {
	path string
}

// NewPinStore returns a PinStore backed by the given file path.
func NewPinStore(path string) *PinStore {
	return &PinStore{path: path}
}

// Pin adds or updates a pinned port entry.
func (s *PinStore) Pin(p PinnedPort) error {
	pins, err := s.Load()
	if err != nil {
		return err
	}
	p.PinnedAt = time.Now().UTC()
	key := pinKey(p.Port, p.Protocol)
	pins[key] = p
	return s.save(pins)
}

// Unpin removes a pinned port entry by port and protocol.
func (s *PinStore) Unpin(port int, protocol string) error {
	pins, err := s.Load()
	if err != nil {
		return err
	}
	delete(pins, pinKey(port, protocol))
	return s.save(pins)
}

// IsPinned returns true if the given port/protocol combination is pinned.
func (s *PinStore) IsPinned(port int, protocol string) (bool, error) {
	pins, err := s.Load()
	if err != nil {
		return false, err
	}
	_, ok := pins[pinKey(port, protocol)]
	return ok, nil
}

// Load reads all pinned ports from disk. Returns an empty map if the file does not exist.
func (s *PinStore) Load() (map[string]PinnedPort, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return make(map[string]PinnedPort), nil
	}
	if err != nil {
		return nil, fmt.Errorf("pin: read %s: %w", s.path, err)
	}
	var pins map[string]PinnedPort
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, fmt.Errorf("pin: unmarshal: %w", err)
	}
	return pins, nil
}

func (s *PinStore) save(pins map[string]PinnedPort) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("pin: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}

func pinKey(port int, protocol string) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
