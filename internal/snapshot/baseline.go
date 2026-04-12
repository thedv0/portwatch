package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline represents a saved reference state of open ports.
type Baseline struct {
	CreatedAt time.Time      `json:"created_at"`
	Label     string         `json:"label"`
	Ports     []scanner.Port `json:"ports"`
}

// BaselineStore manages saving and loading baselines from disk.
type BaselineStore struct {
	dir string
}

// NewBaselineStore returns a BaselineStore rooted at dir.
func NewBaselineStore(dir string) *BaselineStore {
	return &BaselineStore{dir: dir}
}

// Save writes a baseline with the given label to disk.
func (s *BaselineStore) Save(label string, ports []scanner.Port) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("baseline: mkdir: %w", err)
	}
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Ports:     ports,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	path := s.path(label)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load reads a baseline by label from disk.
func (s *BaselineStore) Load(label string) (*Baseline, error) {
	path := s.path(label)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: %q not found", label)
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// List returns all baseline labels available in the store directory.
func (s *BaselineStore) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("baseline: readdir: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			name := e.Name()
			labels = append(labels, name[:len(name)-5])
		}
	}
	return labels, nil
}

func (s *BaselineStore) path(label string) string {
	return filepath.Join(s.dir, label+".json")
}
