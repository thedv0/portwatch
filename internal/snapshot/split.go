package snapshot

import (
	"errors"
	"fmt"
)

// SplitOptions controls how a snapshot slice is divided into partitions.
type SplitOptions struct {
	// NumParts is the number of equal-sized parts to split into.
	NumParts int
	// Field is the field to split on: "protocol", "process", or "port".
	Field string
}

// DefaultSplitOptions returns sensible defaults.
func DefaultSplitOptions() SplitOptions {
	return SplitOptions{
		NumParts: 2,
		Field:    "port",
	}
}

// SplitResult holds one part of the split output.
type SplitResult struct {
	Index int
	Ports []Port
}

// Split divides ports into NumParts roughly equal slices, or groups by Field
// value when Field is "protocol" or "process".
func Split(ports []Port, opts SplitOptions) ([]SplitResult, error) {
	if opts.NumParts < 1 {
		return nil, errors.New("split: NumParts must be >= 1")
	}

	switch opts.Field {
	case "protocol", "process":
		return splitByField(ports, opts.Field)
	case "port", "":
		return splitByCount(ports, opts.NumParts)
	default:
		return nil, fmt.Errorf("split: unknown field %q", opts.Field)
	}
}

func splitByCount(ports []Port, n int) ([]SplitResult, error) {
	results := make([]SplitResult, n)
	for i := range results {
		results[i].Index = i
	}
	for i, p := range ports {
		bucket := i % n
		results[bucket].Ports = append(results[bucket].Ports, p)
	}
	return results, nil
}

func splitByField(ports []Port, field string) ([]SplitResult, error) {
	order := []string{}
	groups := map[string][]Port{}
	for _, p := range ports {
		var key string
		if field == "protocol" {
			key = p.Protocol
		} else {
			key = p.Process
		}
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], p)
	}
	results := make([]SplitResult, len(order))
	for i, k := range order {
		results[i] = SplitResult{Index: i, Ports: groups[k]}
	}
	return results, nil
}
