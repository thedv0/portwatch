package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a key-value label that can be attached to a snapshot entry.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TagSet is an ordered collection of unique tags keyed by name.
type TagSet map[string]string

// NewTagSet parses a slice of "key=value" strings into a TagSet.
// Entries without "=" are stored with an empty value.
func NewTagSet(raw []string) (TagSet, error) {
	ts := make(TagSet, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("tag key must not be empty (got %q)", s)
		}
		val := ""
		if len(parts) == 2 {
			val = strings.TrimSpace(parts[1])
		}
		ts[key] = val
	}
	return ts, nil
}

// Has reports whether the TagSet contains the given key.
func (ts TagSet) Has(key string) bool {
	_, ok := ts[key]
	return ok
}

// Get returns the value for key and whether it was present.
func (ts TagSet) Get(key string) (string, bool) {
	v, ok := ts[key]
	return v, ok
}

// Merge returns a new TagSet combining ts with other; other's values win on conflict.
func (ts TagSet) Merge(other TagSet) TagSet {
	out := make(TagSet, len(ts)+len(other))
	for k, v := range ts {
		out[k] = v
	}
	for k, v := range other {
		out[k] = v
	}
	return out
}

// Slice returns tags as sorted "key=value" strings for deterministic output.
func (ts TagSet) Slice() []string {
	keys := make([]string, 0, len(ts))
	for k := range ts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = k + "=" + ts[k]
	}
	return out
}
