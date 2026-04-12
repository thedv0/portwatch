package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint represents a named point-in-time snapshot of open ports.
type Checkpoint struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Ports     []Port    `json:"ports"`
	Note      string    `json:"note,omitempty"`
}

// CheckpointStore persists named checkpoints to disk.
type CheckpointStore struct {
	dir string
}

// NewCheckpointStore creates a CheckpointStore rooted at dir.
func NewCheckpointStore(dir string) *CheckpointStore {
	return &CheckpointStore{dir: dir}
}

// Save writes a checkpoint to disk under <dir>/<name>.json.
func (s *CheckpointStore) Save(cp Checkpoint) error {
	if cp.Name == "" {
		return fmt.Errorf("checkpoint name must not be empty")
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("create checkpoint dir: %w", err)
	}
	path := filepath.Join(s.dir, cp.Name+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create checkpoint file: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cp)
}

// Load retrieves a checkpoint by name.
func (s *CheckpointStore) Load(name string) (Checkpoint, error) {
	path := filepath.Join(s.dir, name+".json")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Checkpoint{}, fmt.Errorf("checkpoint %q not found", name)
		}
		return Checkpoint{}, fmt.Errorf("open checkpoint: %w", err)
	}
	defer f.Close()
	var cp Checkpoint
	if err := json.NewDecoder(f).Decode(&cp); err != nil {
		return Checkpoint{}, fmt.Errorf("decode checkpoint: %w", err)
	}
	return cp, nil
}

// List returns the names of all stored checkpoints.
func (s *CheckpointStore) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list checkpoints: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a checkpoint by name.
func (s *CheckpointStore) Delete(name string) error {
	path := filepath.Join(s.dir, name+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete checkpoint %q: %w", name, err)
	}
	return nil
}
