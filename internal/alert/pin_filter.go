package alert

import (
	"fmt"

	"github.com/example/portwatch/internal/snapshot"
)

// PinFilter wraps a Matcher and suppresses violations for pinned ports.
type PinFilter struct {
	matcher *Matcher
	store   *snapshot.PinStore
}

// NewPinFilter creates a PinFilter that delegates to the given Matcher and
// skips any violation whose port/protocol pair is pinned in the PinStore.
func NewPinFilter(m *Matcher, store *snapshot.PinStore) *PinFilter {
	return &PinFilter{matcher: m, store: store}
}

// Evaluate runs the underlying matcher and removes violations for pinned ports.
func (f *PinFilter) Evaluate(states []PortState) ([]Event, error) {
	events, err := f.matcher.Evaluate(states)
	if err != nil {
		return nil, err
	}

	filtered := events[:0]
	for _, e := range events {
		pinned, perr := f.store.IsPinned(e.Port, e.Protocol)
		if perr != nil {
			return nil, fmt.Errorf("pin_filter: check pin for %s:%d: %w", e.Protocol, e.Port, perr)
		}
		if !pinned {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}
