package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// AnnotationStore persists port annotations to a JSON file.
type AnnotationStore struct {
	mu   sync.RWMutex
	path string
	data map[string][]Annotation // keyed by portKey (e.g. "tcp:8080")
}

// NewAnnotationStore creates a store backed by the given file path.
func NewAnnotationStore(path string) *AnnotationStore {
	return &AnnotationStore{path: path, data: make(map[string][]Annotation)}
}

// Load reads persisted annotations from disk; missing file is not an error.
func (s *AnnotationStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("annotate_store: read %s: %w", s.path, err)
	}
	return json.Unmarshal(b, &s.data)
}

// Save writes all annotations to disk atomically.
func (s *AnnotationStore) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("annotate_store: mkdir: %w", err)
	}
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

// Set replaces the annotation set for the given protocol and port.
func (s *AnnotationStore) Set(proto string, port int, as *AnnotationSet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := fmt.Sprintf("%s:%d", proto, port)
	if as == nil || as.Len() == 0 {
		delete(s.data, key)
		return
	}
	s.data[key] = as.All()
}

// Get retrieves the AnnotationSet for a given protocol and port.
func (s *AnnotationStore) Get(proto string, port int) *AnnotationSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := fmt.Sprintf("%s:%d", proto, port)
	anns, ok := s.data[key]
	if !ok {
		return &AnnotationSet{}
	}
	as := &AnnotationSet{annotations: make([]Annotation, len(anns))}
	copy(as.annotations, anns)
	return as
}

// Keys returns all stored port keys.
func (s *AnnotationStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}
