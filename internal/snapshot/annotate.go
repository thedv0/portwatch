package snapshot

import (
	"fmt"
	"strings"
)

// Annotation holds a key-value note attached to a port entry.
type Annotation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AnnotationSet is an ordered collection of annotations for a port.
type AnnotationSet struct {
	annotations []Annotation
}

// NewAnnotationSet parses annotations from "key=value" strings.
func NewAnnotationSet(pairs []string) (*AnnotationSet, error) {
	as := &AnnotationSet{}
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid annotation %q: expected key=value", p)
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("annotation key must not be empty in %q", p)
		}
		as.annotations = append(as.annotations, Annotation{Key: key, Value: strings.TrimSpace(parts[1])})
	}
	return as, nil
}

// Add appends an annotation, overwriting any existing entry with the same key.
func (as *AnnotationSet) Add(key, value string) {
	for i, a := range as.annotations {
		if a.Key == key {
			as.annotations[i].Value = value
			return
		}
	}
	as.annotations = append(as.annotations, Annotation{Key: key, Value: value})
}

// Get returns the value for a key and whether it was found.
func (as *AnnotationSet) Get(key string) (string, bool) {
	for _, a := range as.annotations {
		if a.Key == key {
			return a.Value, true
		}
	}
	return "", false
}

// All returns a copy of all annotations.
func (as *AnnotationSet) All() []Annotation {
	out := make([]Annotation, len(as.annotations))
	copy(out, as.annotations)
	return out
}

// Len returns the number of annotations.
func (as *AnnotationSet) Len() int {
	return len(as.annotations)
}

// Remove deletes an annotation by key, returning true if it existed.
func (as *AnnotationSet) Remove(key string) bool {
	for i, a := range as.annotations {
		if a.Key == key {
			as.annotations = append(as.annotations[:i], as.annotations[i+1:]...)
			return true
		}
	}
	return false
}
