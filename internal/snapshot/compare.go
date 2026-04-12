package snapshot

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// CompareResult holds the outcome of comparing two port snapshots.
type CompareResult struct {
	Added   []scanner.Port
	Removed []scanner.Port
	Changed []PortChange
}

// PortChange describes a port whose attributes changed between snapshots.
type PortChange struct {
	Port    scanner.Port
	OldPID  int
	NewPID  int
	OldProc string
	NewProc string
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	var sb strings.Builder
	if len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0 {
		sb.WriteString("no changes detected")
		return sb.String()
	}
	if len(r.Added) > 0 {
		fmt.Fprintf(&sb, "added %d port(s)\n", len(r.Added))
		for _, p := range r.Added {
			fmt.Fprintf(&sb, "  + %s/%d (pid=%d proc=%s)\n", p.Protocol, p.Port, p.PID, p.Process)
		}
	}
	if len(r.Removed) > 0 {
		fmt.Fprintf(&sb, "removed %d port(s)\n", len(r.Removed))
		for _, p := range r.Removed {
			fmt.Fprintf(&sb, "  - %s/%d (pid=%d proc=%s)\n", p.Protocol, p.Port, p.PID, p.Process)
		}
	}
	if len(r.Changed) > 0 {
		fmt.Fprintf(&sb, "changed %d port(s)\n", len(r.Changed))
		for _, c := range r.Changed {
			fmt.Fprintf(&sb, "  ~ %s/%d pid:%d->%d proc:%s->%s\n",
				c.Port.Protocol, c.Port.Port, c.OldPID, c.NewPID, c.OldProc, c.NewProc)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Compare performs a detailed comparison between two port slices, tracking
// additions, removals, and attribute changes (PID/process) for the same port key.
func Compare(prev, curr []scanner.Port) CompareResult {
	type meta struct {
		pid  int
		proc string
	}

	prevMap := make(map[string]meta, len(prev))
	for _, p := range prev {
		prevMap[portKey(p)] = meta{pid: p.PID, proc: p.Process}
	}

	currMap := make(map[string]meta, len(curr))
	for _, p := range curr {
		currMap[portKey(p)] = meta{pid: p.PID, proc: p.Process}
	}

	var result CompareResult

	for _, p := range curr {
		key := portKey(p)
		old, exists := prevMap[key]
		if !exists {
			result.Added = append(result.Added, p)
			continue
		}
		if old.pid != p.PID || old.proc != p.Process {
			result.Changed = append(result.Changed, PortChange{
				Port:    p,
				OldPID:  old.pid,
				NewPID:  p.PID,
				OldProc: old.proc,
				NewProc: p.Process,
			})
		}
	}

	for _, p := range prev {
		if _, exists := currMap[portKey(p)]; !exists {
			result.Removed = append(result.Removed, p)
		}
	}

	return result
}
