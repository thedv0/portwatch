package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeSumPorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Protocol: "tcp", PID: 100, Process: "nginx"},
		{Port: 443, Protocol: "tcp", PID: 100, Process: "nginx"},
		{Port: 53, Protocol: "udp", PID: 200, Process: "dnsmasq"},
		{Port: 53, Protocol: "tcp", PID: 200, Process: "dnsmasq"},
		{Port: 8080, Protocol: "tcp", PID: 300, Process: ""},
	}
}

func TestSummarize_TotalPorts(t *testing.T) {
	ports := makeSumPorts()
	s := Summarize(ports)
	if s.TotalPorts != 5 {
		t.Errorf("expected 5 total ports, got %d", s.TotalPorts)
	}
}

func TestSummarize_ByProtocol(t *testing.T) {
	s := Summarize(makeSumPorts())
	if s.ByProtocol["tcp"] != 4 {
		t.Errorf("expected 4 tcp, got %d", s.ByProtocol["tcp"])
	}
	if s.ByProtocol["udp"] != 1 {
		t.Errorf("expected 1 udp, got %d", s.ByProtocol["udp"])
	}
}

func TestSummarize_ByProcess(t *testing.T) {
	s := Summarize(makeSumPorts())
	if s.ByProcess["nginx"] != 2 {
		t.Errorf("expected nginx=2, got %d", s.ByProcess["nginx"])
	}
	if s.ByProcess["dnsmasq"] != 2 {
		t.Errorf("expected dnsmasq=2, got %d", s.ByProcess["dnsmasq"])
	}
	if s.ByProcess["pid:300"] != 1 {
		t.Errorf("expected pid:300=1, got %d", s.ByProcess["pid:300"])
	}
}

func TestSummarize_TopPorts_OrderedByFrequency(t *testing.T) {
	s := Summarize(makeSumPorts())
	if len(s.TopPorts) == 0 {
		t.Fatal("expected non-empty TopPorts")
	}
	// port 53 appears twice, should be first
	if s.TopPorts[0].Port != 53 {
		t.Errorf("expected port 53 first, got %d", s.TopPorts[0].Port)
	}
	if s.TopPorts[0].Count != 2 {
		t.Errorf("expected count 2, got %d", s.TopPorts[0].Count)
	}
}

func TestSummarize_ListeningPIDs_Unique(t *testing.T) {
	s := Summarize(makeSumPorts())
	if len(s.ListeningPIDs) != 3 {
		t.Errorf("expected 3 unique PIDs, got %d", len(s.ListeningPIDs))
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.TotalPorts != 0 {
		t.Errorf("expected 0 total ports, got %d", s.TotalPorts)
	}
	if len(s.TopPorts) != 0 {
		t.Errorf("expected empty TopPorts")
	}
}
