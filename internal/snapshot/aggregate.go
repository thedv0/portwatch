package snapshot

import "time"

// AggregateStats holds rolled-up statistics across multiple snapshots.
type AggregateStats struct {
	From        time.Time
	To          time.Time
	SampleCount int
	AvgOpen     float64
	MaxOpen     int
	MinOpen     int
	TopPorts    []PortFrequency
	TopProcesses []ProcessFrequency
}

// PortFrequency pairs a port number with how often it appeared.
type PortFrequency struct {
	Port  int
	Count int
}

// ProcessFrequency pairs a process name with how often it appeared.
type ProcessFrequency struct {
	Process string
	Count   int
}

// Aggregate computes rolled-up statistics across a slice of historical snapshots.
func Aggregate(snaps []Snapshot) AggregateStats {
	if len(snaps) == 0 {
		return AggregateStats{}
	}

	portCounts := make(map[int]int)
	procCounts := make(map[string]int)
	totalOpen := 0
	maxOpen := 0
	minOpen := -1

	from := snaps[0].Timestamp
	to := snaps[0].Timestamp

	for _, s := range snaps {
		if s.Timestamp.Before(from) {
			from = s.Timestamp
		}
		if s.Timestamp.After(to) {
			to = s.Timestamp
		}
		count := len(s.Ports)
		totalOpen += count
		if count > maxOpen {
			maxOpen = count
		}
		if minOpen < 0 || count < minOpen {
			minOpen = count
		}
		for _, p := range s.Ports {
			portCounts[p.Port]++
			if p.Process != "" {
				procCounts[p.Process]++
			}
		}
	}

	if minOpen < 0 {
		minOpen = 0
	}

	return AggregateStats{
		From:         from,
		To:           to,
		SampleCount:  len(snaps),
		AvgOpen:      float64(totalOpen) / float64(len(snaps)),
		MaxOpen:      maxOpen,
		MinOpen:      minOpen,
		TopPorts:     topPorts(portCounts, 5),
		TopProcesses: topProcesses(procCounts, 5),
	}
}

func topPorts(counts map[int]int, n int) []PortFrequency {
	result := make([]PortFrequency, 0, len(counts))
	for port, c := range counts {
		result = append(result, PortFrequency{Port: port, Count: c})
	}
	sortPortFreq(result)
	if len(result) > n {
		result = result[:n]
	}
	return result
}

func topProcesses(counts map[string]int, n int) []ProcessFrequency {
	result := make([]ProcessFrequency, 0, len(counts))
	for proc, c := range counts {
		result = append(result, ProcessFrequency{Process: proc, Count: c})
	}
	sortProcFreq(result)
	if len(result) > n {
		result = result[:n]
	}
	return result
}

func sortPortFreq(s []PortFrequency) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j].Count > s[j-1].Count; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func sortProcFreq(s []ProcessFrequency) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j].Count > s[j-1].Count; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
