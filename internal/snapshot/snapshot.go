package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// State represents a saved snapshot of open ports at a point in time.
type State struct {
	Timestamp time.Time        `json:"timestamp"`
	Ports     []scanner.Port   `json:"ports"`
}

// Store handles persistence of port snapshots to disk.
type Store struct {
	path string
}

// NewStore creates a Store that persists snapshots at the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the current port list to disk as a JSON snapshot.
func (s *Store) Save(ports []scanner.Port) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	state := State{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
	f, err := os.CreateTemp(filepath.Dir(s.path), ".portwatch-snap-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := json.NewEncoder(f).Encode(state); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, s.path)
}

// Load reads the last saved snapshot from disk.
// Returns an empty State (zero value) when no snapshot exists yet.
func (s *Store) Load() (State, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, err
	}
	return state, nil
}
