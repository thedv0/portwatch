package snapshot

import "fmt"

// CheckpointDiffResult holds the comparison between a checkpoint and current ports.
type CheckpointDiffResult struct {
	CheckpointName string
	Added          []Port // present now, absent in checkpoint
	Removed        []Port // present in checkpoint, absent now
	Unchanged      int
}

// CompareToCheckpoint diffs current ports against a stored checkpoint.
func CompareToCheckpoint(store *CheckpointStore, name string, current []Port) (CheckpointDiffResult, error) {
	cp, err := store.Load(name)
	if err != nil {
		return CheckpointDiffResult{}, fmt.Errorf("load checkpoint for diff: %w", err)
	}

	result := CheckpointDiffResult{CheckpointName: name}

	baseline := indexPorts(cp.Ports)
	now := indexPorts(current)

	for k, p := range now {
		if _, ok := baseline[k]; !ok {
			result.Added = append(result.Added, p)
		} else {
			result.Unchanged++
		}
	}
	for k, p := range baseline {
		if _, ok := now[k]; !ok {
			result.Removed = append(result.Removed, p)
		}
	}
	return result, nil
}
