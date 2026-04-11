package snapshot

import (
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Summary holds aggregated statistics for a port snapshot.
type Summary struct {
	TotalPorts   int
	ByProtocol   map[string]int
	ByProcess    map[string]int
	TopPorts     []PortCount
	ListeningPIDs []int
}

// PortCount pairs a port number with how many times it appears.
type PortCount struct {
	Port  int
	Count int
}

// Summarize computes aggregate statistics from a slice of ports.
func Summarize(ports []scanner.Port) Summary {
	s := Summary{
		TotalPorts: len(ports),
		ByProtocol: make(map[string]int),
		ByProcess:  make(map[string]int),
	}

	portFreq := make(map[int]int)
	pidSeen := make(map[int]bool)

	for _, p := range ports {
		s.ByProtocol[p.Protocol]++

		name := p.Process
		if name == "" {
			name = fmt.Sprintf("pid:%d", p.PID)
		}
		s.ByProcess[name]++

		portFreq[p.Port]++

		if p.PID > 0 && !pidSeen[p.PID] {
			pidSeen[p.PID] = true
			s.ListeningPIDs = append(s.ListeningPIDs, p.PID)
		}
	}

	for port, count := range portFreq {
		s.TopPorts = append(s.TopPorts, PortCount{Port: port, Count: count})
	}
	sort.Slice(s.TopPorts, func(i, j int) bool {
		if s.TopPorts[i].Count != s.TopPorts[j].Count {
			return s.TopPorts[i].Count > s.TopPorts[j].Count
		}
		return s.TopPorts[i].Port < s.TopPorts[j].Port
	})
	if len(s.TopPorts) > 10 {
		s.TopPorts = s.TopPorts[:10]
	}

	sort.Ints(s.ListeningPIDs)
	return s
}
