package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// PortEntry represents a single open port detected on the system.
type PortEntry struct {
	Protocol string
	LocalAddr string
	Port      int
	PID       int
	Process   string
}

// String returns a human-readable representation of a PortEntry.
func (p PortEntry) String() string {
	return fmt.Sprintf("%s %s:%d (pid=%d, process=%s)", p.Protocol, p.LocalAddr, p.Port, p.PID, p.Process)
}

// Scanner scans for open listening ports on the local machine.
type Scanner struct {
	Timeout time.Duration
}

// NewScanner creates a Scanner with a default timeout.
func NewScanner(timeout time.Duration) *Scanner {
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	return &Scanner{Timeout: timeout}
}

// ScanTCP probes a range of TCP ports and returns those that are open.
// In production this would read from /proc/net/tcp; here we use dial-based detection.
func (s *Scanner) ScanTCP(ports []int) ([]PortEntry, error) {
	var entries []PortEntry
	for _, port := range ports {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err != nil {
			continue
		}
		conn.Close()
		host, portStr, _ := net.SplitHostPort(addr)
		p, _ := strconv.Atoi(portStr)
		entries = append(entries, PortEntry{
			Protocol:  "tcp",
			LocalAddr: host,
			Port:      p,
			PID:       -1,
			Process:   "unknown",
		})
	}
	return entries, nil
}

// ParsePortRange parses a string like "80,443,8000-8080" into a slice of port numbers.
func ParsePortRange(raw string) ([]int, error) {
	var ports []int
	seen := map[int]bool{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err1 := strconv.Atoi(bounds[0])
			hi, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || lo > hi {
				return nil, fmt.Errorf("invalid range: %q", part)
			}
			for i := lo; i <= hi; i++ {
				if !seen[i] {
					ports = append(ports, i)
					seen[i] = true
				}
			}
		} else {
			n, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %q", part)
			}
			if !seen[n] {
				ports = append(ports, n)
				seen[n] = true
			}
		}
	}
	return ports, nil
}
