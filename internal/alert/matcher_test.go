package alert_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/rules"
)

func makeConfig(allowedPorts []int, proto string) *rules.Config {
	return &rules.Config{
		Rules: []rules.Rule{
			{
				Name:     "allowed",
				Ports:    allowedPorts,
				Protocol: proto,
				Alert:    false,
			},
		},
	}
}

func TestEvaluate_NoViolations(t *testing.T) {
	cfg := makeConfig([]int{22, 80, 443}, "tcp")
	m := alert.NewMatcher(cfg)

	states := []alert.PortState{
		{Port: 22, Protocol: "tcp", Process: "sshd"},
		{Port: 80, Protocol: "tcp", Process: "nginx"},
	}

	events := m.Evaluate(states)
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d: %v", len(events), events)
	}
}

func TestEvaluate_UnexpectedPort(t *testing.T) {
	cfg := makeConfig([]int{22, 80}, "tcp")
	m := alert.NewMatcher(cfg)

	states := []alert.PortState{
		{Port: 22, Protocol: "tcp", Process: "sshd"},
		{Port: 4444, Protocol: "tcp", Process: "nc"},
	}

	events := m.Evaluate(states)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Port != 4444 {
		t.Errorf("expected port 4444, got %d", events[0].Port)
	}
	if events[0].Level != alert.LevelAlert {
		t.Errorf("expected ALERT level, got %s", events[0].Level)
	}
}

func TestEvaluate_EmptyStates(t *testing.T) {
	cfg := makeConfig([]int{22}, "tcp")
	m := alert.NewMatcher(cfg)

	events := m.Evaluate(nil)
	if len(events) != 0 {
		t.Errorf("expected 0 events for nil states, got %d", len(events))
	}
}

func TestEvaluate_ProtocolMismatch(t *testing.T) {
	cfg := makeConfig([]int{53}, "tcp")
	m := alert.NewMatcher(cfg)

	// Port 53 allowed on tcp but detected on udp — should alert.
	states := []alert.PortState{
		{Port: 53, Protocol: "udp", Process: "dnsmasq"},
	}

	events := m.Evaluate(states)
	if len(events) != 1 {
		t.Errorf("expected 1 event for protocol mismatch, got %d", len(events))
	}
}
