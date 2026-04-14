package snapshot

import "time"

// CapOptions controls how a port list is capped before storage or reporting.
type CapOptions struct {
	// MaxPorts is the maximum number of ports to retain. Zero means no limit.
	MaxPorts int
	// SortByPort sorts the list by port number before capping so the lowest
	// numbered ports are always retained.
	SortByPort bool
	// PreferLowPID retains entries with lower PIDs when deduplicating before cap.
	PreferLowPID bool
}

// DefaultCapOptions returns a sensible default cap configuration.
func DefaultCapOptions() CapOptions {
	return CapOptions{
		MaxPorts:   0,
		SortByPort: true,
	}
}

// Cap limits the number of ports in a slice according to opts.
// When SortByPort is true the slice is sorted ascending by port number before
// truncation so the lowest ports survive the cap.
func Cap(ports []PortState, opts CapOptions) []PortState {
	if len(ports) == 0 {
		return ports
	}

	out := make([]PortState, len(ports))
	copy(out, ports)

	if opts.SortByPort {
		sortPortsByNumber(out)
	}

	if opts.MaxPorts > 0 && len(out) > opts.MaxPorts {
		out = out[:opts.MaxPorts]
	}

	return out
}

// sortPortsByNumber performs a simple insertion sort on port number.
func sortPortsByNumber(ports []PortState) {
	for i := 1; i < len(ports); i++ {
		for j := i; j > 0 && ports[j].Port < ports[j-1].Port; j-- {
			ports[j], ports[j-1] = ports[j-1], ports[j]
		}
	}
}

// CapSnapshot applies Cap to a snapshot's port list and returns a new snapshot.
func CapSnapshot(snap Snapshot, opts CapOptions) Snapshot {
	return Snapshot{
		Timestamp: snap.Timestamp,
		Ports:     Cap(snap.Ports, opts),
	}
}

// Snapshot is a lightweight local alias used within this file to avoid import
// cycles; the real type is defined in snapshot.go.
type capSnapshotRef = struct {
	Timestamp time.Time
	Ports     []PortState
}
