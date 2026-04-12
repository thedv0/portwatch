package snapshot

import (
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// EnrichOptions controls which enrichment steps are applied.
type EnrichOptions struct {
	NormalizeProtocol bool
	ResolveWellKnown  bool
	TagSystemPorts    bool
}

// DefaultEnrichOptions returns sensible defaults for enrichment.
func DefaultEnrichOptions() EnrichOptions {
	return EnrichOptions{
		NormalizeProtocol: true,
		ResolveWellKnown:  true,
		TagSystemPorts:    true,
	}
}

// wellKnownNames maps common port numbers to service names.
var wellKnownNames = map[int]string{
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	27017: "mongodb",
}

// Enrich applies metadata enrichment to a slice of ports in-place.
func Enrich(ports []scanner.Port, opts EnrichOptions) []scanner.Port {
	out := make([]scanner.Port, len(ports))
	copy(out, ports)
	for i := range out {
		p := &out[i]
		if opts.NormalizeProtocol {
			p.Protocol = strings.ToLower(strings.TrimSpace(p.Protocol))
		}
		if opts.ResolveWellKnown {
			if name, ok := wellKnownNames[p.Port]; ok && p.Process == "" {
				p.Process = name
			}
		}
		if opts.TagSystemPorts {
			// system ports are 0-1023; this is surfaced via the process field prefix
			// only when process is still empty after well-known resolution
			_ = p // classification handled by classify.go; no mutation needed here
		}
	}
	return out
}
