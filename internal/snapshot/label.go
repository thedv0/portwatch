package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// Label represents a key-value annotation attached to a port snapshot entry.
type Label struct {
	Key   string
	Value string
}

// LabelSet holds an ordered, deduplicated collection of labels.
type LabelSet struct {
	labels map[string]string
}

// NewLabelSet creates a LabelSet from a slice of "key=value" strings.
// Keys must be non-empty; values may be empty.
func NewLabelSet(pairs []string) (*LabelSet, error) {
	ls := &LabelSet{labels: make(map[string]string)}
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid label %q: must be key=value", p)
		}
		ls.labels[parts[0]] = parts[1]
	}
	return ls, nil
}

// Set adds or overwrites a label.
func (ls *LabelSet) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("label key must not be empty")
	}
	ls.labels[key] = value
	return nil
}

// Get returns the value for key and whether it was found.
func (ls *LabelSet) Get(key string) (string, bool) {
	v, ok := ls.labels[key]
	return v, ok
}

// Delete removes a label by key.
func (ls *LabelSet) Delete(key string) {
	delete(ls.labels, key)
}

// All returns all labels sorted by key.
func (ls *LabelSet) All() []Label {
	keys := make([]string, 0, len(ls.labels))
	for k := range ls.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]Label, 0, len(keys))
	for _, k := range keys {
		out = append(out, Label{Key: k, Value: ls.labels[k]})
	}
	return out
}

// Merge copies labels from other into ls, overwriting duplicates.
func (ls *LabelSet) Merge(other *LabelSet) {
	for k, v := range other.labels {
		ls.labels[k] = v
	}
}

// Len returns the number of labels.
func (ls *LabelSet) Len() int {
	return len(ls.labels)
}
