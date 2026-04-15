package snapshot

import (
	"github.com/user/portwatch/internal/scanner"
)

// IntersectOptions controls how snapshot intersection is performed.
type IntersectOptions struct {
	// KeyFields determines which port fields form the match key.
	// Supported values: "port", "protocol", "pid", "process"
	KeyFields []string
}

// DefaultIntersectOptions returns sensible defaults: match by port and protocol.
func DefaultIntersectOptions() IntersectOptions {
	return IntersectOptions{
		KeyFields: []string{"port", "protocol"},
	}
}

// intersectKey builds a lookup key from the selected fields.
func intersectKey(p scanner.Port, fields []string) string {
	key := ""
	for _, f := range fields {
		switch f {
		case "port":
			key += itoa(p.Port) + "|"
		case "protocol":
			key += p.Protocol + "|"
		case "pid":
			key += itoa(p.PID) + "|"
		case "process":
			key += p.Process + "|"
		}
	}
	return key
}

// Intersect returns only the ports that appear in every provided snapshot,
// matched according to opts.KeyFields. The returned slice preserves the
// order from the first snapshot.
func Intersect(snaps []Snapshot, opts IntersectOptions) []scanner.Port {
	if len(snaps) == 0 {
		return nil
	}
	if len(opts.KeyFields) == 0 {
		opts = DefaultIntersectOptions()
	}

	// Index keys present in each snapshot after the first.
	presence := make([]map[string]struct{}, len(snaps)-1)
	for i := 1; i < len(snaps); i++ {
		m := make(map[string]struct{}, len(snaps[i].Ports))
		for _, p := range snaps[i].Ports {
			m[intersectKey(p, opts.KeyFields)] = struct{}{}
		}
		presence[i-1] = m
	}

	var result []scanner.Port
	for _, p := range snaps[0].Ports {
		k := intersectKey(p, opts.KeyFields)
		shared := true
		for _, m := range presence {
			if _, ok := m[k]; !ok {
				shared = false
				break
			}
		}
		if shared {
			result = append(result, p)
		}
	}
	return result
}
