package snapshot

import "github.com/user/portwatch/internal/scanner"

// RiskClass represents a broad risk classification for a port.
type RiskClass string

const (
	RiskClassSystem    RiskClass = "system"    // ports 1-1023
	RiskClassRegistered RiskClass = "registered" // ports 1024-49151
	RiskClassDynamic   RiskClass = "dynamic"   // ports 49152-65535
	RiskClassUnknown   RiskClass = "unknown"
)

// ClassifyOptions controls classification behaviour.
type ClassifyOptions struct {
	// WellKnownPorts is an optional allow-list of ports considered "expected".
	WellKnownPorts []int
}

// ClassifiedPort pairs a scanned port with its classification.
type ClassifiedPort struct {
	Port      scanner.Port
	Class     RiskClass
	WellKnown bool
}

// DefaultClassifyOptions returns sensible defaults.
func DefaultClassifyOptions() ClassifyOptions {
	return ClassifyOptions{
		WellKnownPorts: []int{22, 80, 443, 8080, 8443},
	}
}

// Classify assigns a RiskClass and WellKnown flag to each port.
func Classify(ports []scanner.Port, opts ClassifyOptions) []ClassifiedPort {
	wellKnown := make(map[int]bool, len(opts.WellKnownPorts))
	for _, p := range opts.WellKnownPorts {
		wellKnown[p] = true
	}

	out := make([]ClassifiedPort, 0, len(ports))
	for _, p := range ports {
		out = append(out, ClassifiedPort{
			Port:      p,
			Class:     classifyPort(p.Port),
			WellKnown: wellKnown[p.Port],
		})
	}
	return out
}

func classifyPort(port int) RiskClass {
	switch {
	case port >= 1 && port <= 1023:
		return RiskClassSystem
	case port >= 1024 && port <= 49151:
		return RiskClassRegistered
	case port >= 49152 && port <= 65535:
		return RiskClassDynamic
	default:
		return RiskClassUnknown
	}
}
