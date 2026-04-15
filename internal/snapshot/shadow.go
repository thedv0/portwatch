package snapshot

import (
	"fmt"
	"time"
)

// ShadowOptions controls how shadow comparisons are performed.
type ShadowOptions struct {
	// Tolerance is the number of ports that may differ before flagging divergence.
	Tolerance int
	// IgnoreProcess skips process name differences when comparing.
	IgnoreProcess bool
	// Clock is used to stamp the result; defaults to time.Now.
	Clock func() time.Time
}

// DefaultShadowOptions returns sensible defaults.
func DefaultShadowOptions() ShadowOptions {
	return ShadowOptions{
		Tolerance:     0,
		IgnoreProcess: false,
		Clock:         time.Now,
	}
}

// ShadowResult holds the outcome of a shadow comparison.
type ShadowResult struct {
	Timestamp   time.Time
	Diverged    bool
	Divergences []ShadowDivergence
	PrimaryLen  int
	ShadowLen   int
}

// ShadowDivergence describes a single port-level difference.
type ShadowDivergence struct {
	Port     int
	Proteason   string
}

// Shadow compares a primary snapshot against a shadow snapshot and reports
// divergences. It is useful for validating a new scanner implementation
// against a known-good one.
func Shadow(primary, shadow []PortState, opts ShadowOptions) (ShadowResult, error) {
	if opts.Clock == nil {
		return ShadowResult{}, fmt.Errorf("shadow: Clock must not be nil")
	}

	pIdx := make(map[string]PortState, len(primary))
	for _, p := range primary {
		pIdx[shadowKey(p)] = p
	}
	sIdx := make(map[string]PortState, len(shadow))
	for _, s := range shadow {
		sIdx[shadowKey(s)] = s
	}

	var divs []ShadowDivergence

	for k, sp := range sIdx {
		if _, ok := pIdx[k]; !ok {
			divs = append(divs, ShadowDivergence{
				Port:     sp.Port,
				Protocol: sp.Protocol,
				Reason:   "present in shadow, missing in primary",
			})
		}
	}

	for k, pp := range pIdx {
		if ss, ok := sIdx[k]; !ok {
			divs = append(divs, ShadowDivergence{
				Port:     pp.Port,
				Protocol: pp.Protocol,
				Reason:   "present in primary, missing in shadow",
			})
		} else if !opts.IgnoreProcess && pp.Process != ss.Process {
			divs = append(divs, ShadowDivergence{
				Port:     pp.Port,
				Protocol: pp.Protocol,
				Reason:   fmt.Sprintf("process mismatch: primary=%q shadow=%q", pp.Process, ss.Process),
			})
		}
	}

	return ShadowResult{
		Timestamp:   opts.Clock(),
		Diverged:    len(divs) > opts.Tolerance,
		Divergences: divs,
		PrimaryLen:  len(primary),
		ShadowLen:   len(shadow),
	}, nil
}

func shadowKey(p PortState) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Port)
}
