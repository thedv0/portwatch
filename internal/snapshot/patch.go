package snapshot

import (
	"fmt"
	"time"
)

// PatchOp represents a single mutation operation on a port entry.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
)

// Patch describes a targeted mutation for a specific port key.
type Patch struct {
	Op      PatchOp
	Key     string // e.g. "process", "pid", "protocol"
	Value   string
	PortNum int
	Proto   string
}

// PatchResult holds the outcome of applying patches to a snapshot.
type PatchResult struct {
	Timestamp time.Time
	Applied   int
	Skipped   int
	Ports     []PortState
}

// DefaultPatchOptions returns a PatchOptions with safe defaults.
type PatchOptions struct {
	IgnoreMissing bool
}

func DefaultPatchOptions() PatchOptions {
	return PatchOptions{IgnoreMissing: true}
}

// ApplyPatches applies a list of Patch operations to the given ports.
// Ports are matched by (PortNum, Proto). Unknown keys are skipped.
func ApplyPatches(ports []PortState, patches []Patch, opts PatchOptions) (PatchResult, error) {
	result := PatchResult{
		Timestamp: time.Now(),
	}

	copied := make([]PortState, len(ports))
	copy(copied, ports)

	index := make(map[string]int, len(copied))
	for i, p := range copied {
		index[patchKey(p.Port, p.Protocol)] = i
	}

	for _, patch := range patches {
		k := patchKey(patch.PortNum, patch.Proto)
		idx, ok := index[k]
		if !ok {
			if opts.IgnoreMissing {
				result.Skipped++
				continue
			}
			return result, fmt.Errorf("patch target not found: port %d/%s", patch.PortNum, patch.Proto)
		}

		switch patch.Op {
		case PatchSet:
			switch patch.Key {
			case "process":
				copied[idx].Process = patch.Value
			case "protocol":
				copied[idx].Protocol = patch.Value
			default:
				result.Skipped++
				continue
			}
		case PatchDelete:
			switch patch.Key {
			case "process":
				copied[idx].Process = ""
			default:
				result.Skipped++
				continue
			}
		default:
			result.Skipped++
			continue
		}
		result.Applied++
	}

	result.Ports = copied
	return result, nil
}

func patchKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
