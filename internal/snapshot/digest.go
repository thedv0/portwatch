package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

// DigestOptions controls how a snapshot digest is computed.
type DigestOptions struct {
	// IncludePID includes PID values in the digest computation.
	IncludePID bool
	// IncludeProcess includes process names in the digest computation.
	IncludeProcess bool
}

// DefaultDigestOptions returns sensible defaults.
func DefaultDigestOptions() DigestOptions {
	return DigestOptions{
		IncludePID:     false,
		IncludeProcess: true,
	}
}

// DigestResult holds the output of a Digest operation.
type DigestResult struct {
	Hex      string
	PortCount int
}

// Digest computes a stable SHA-256 digest over a slice of PortState values.
// Ports are sorted before hashing to ensure deterministic output regardless
// of input order.
func Digest(ports []PortState, opts DigestOptions) DigestResult {
	if len(ports) == 0 {
		return DigestResult{Hex: emptyDigest()}
	}

	sorted := make([]PortState, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Protocol != sorted[j].Protocol {
			return sorted[i].Protocol < sorted[j].Protocol
		}
		return sorted[i].Port < sorted[j].Port
	})

	h := sha256.New()
	for _, p := range sorted {
		parts := p.Protocol + ":" + strconv.Itoa(p.Port)
		if opts.IncludeProcess {
			parts += ":" + p.Process
		}
		if opts.IncludePID {
			parts += ":" + strconv.Itoa(p.PID)
		}
		fmt.Fprintln(h, parts)
	}

	return DigestResult{
		Hex:      hex.EncodeToString(h.Sum(nil)),
		PortCount: len(sorted),
	}
}

func emptyDigest() string {
	h := sha256.New()
	return hex.EncodeToString(h.Sum(nil))
}
