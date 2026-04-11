package snapshot

import (
	"testing"
)

func makeQueryPorts() []PortState {
	return []PortState{
		{Port: 443, Protocol: "tcp", PID: 200, Process: "nginx"},
		{Port: 80, Protocol: "tcp", PID: 100, Process: "apache"},
		{Port: 8080, Protocol: "udp", PID: 300, Process: "myapp"},
		{Port: 22, Protocol: "tcp", PID: 50, Process: "sshd"},
	}
}

func TestQuery_SortByPort_Ascending(t *testing.T) {
	ports := makeQueryPorts()
	opts := DefaultQueryOptions()
	result := Query(ports, opts)
	if result[0].Port != 22 || result[1].Port != 80 || result[2].Port != 443 || result[3].Port != 8080 {
		t.Errorf("unexpected port order: %v", result)
	}
}

func TestQuery_SortByPort_Descending(t *testing.T) {
	ports := makeQueryPorts()
	opts := DefaultQueryOptions()
	opts.Ascending = false
	result := Query(ports, opts)
	if result[0].Port != 8080 || result[3].Port != 22 {
		t.Errorf("unexpected descending order: %v", result)
	}
}

func TestQuery_SortByPID(t *testing.T) {
	ports := makeQueryPorts()
	opts := DefaultQueryOptions()
	opts.SortBy = SortByPID
	result := Query(ports, opts)
	if result[0].PID != 50 || result[3].PID != 300 {
		t.Errorf("unexpected PID order: %v", result)
	}
}

func TestQuery_SortByProcess(t *testing.T) {
	ports := makeQueryPorts()
	opts := DefaultQueryOptions()
	opts.SortBy = SortByProcess
	result := Query(ports, opts)
	if result[0].Process != "apache" {
		t.Errorf("expected apache first, got %s", result[0].Process)
	}
}

func TestQuery_LimitAndOffset(t *testing.T) {
	ports := makeQueryPorts()
	opts := DefaultQueryOptions()
	opts.Limit = 2
	opts.Offset = 1
	result := Query(ports, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80 at index 0, got %d", result[0].Port)
	}
}

func TestQuery_OffsetBeyondLength(t *testing.T) {
	ports := makeQueryPorts()
	opts := DefaultQueryOptions()
	opts.Offset = 100
	result := Query(ports, opts)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestQuery_EmptyInput(t *testing.T) {
	result := Query([]PortState{}, DefaultQueryOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input")
	}
}

func TestDefaultQueryOptions_Values(t *testing.T) {
	opts := DefaultQueryOptions()
	if opts.SortBy != SortByPort {
		t.Errorf("expected SortByPort, got %s", opts.SortBy)
	}
	if !opts.Ascending {
		t.Error("expected ascending by default")
	}
	if opts.Limit != 0 || opts.Offset != 0 {
		t.Error("expected zero limit and offset by default")
	}
}
