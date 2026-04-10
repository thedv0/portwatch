package alert

import (
	"github.com/yourorg/portwatch/internal/rules"
)

// PortState holds the current state of a single open port.
type PortState struct {
	Port     int
	Protocol string
	Process  string
}

// Matcher evaluates detected ports against a loaded configuration
// and returns alert Events for any violations.
type Matcher struct {
	cfg *rules.Config
}

// NewMatcher creates a Matcher backed by the given Config.
func NewMatcher(cfg *rules.Config) *Matcher {
	return &Matcher{cfg: cfg}
}

// Evaluate checks a slice of PortState values against the allow-list rules
// and returns Events for each unexpected listener.
func (m *Matcher) Evaluate(states []PortState) []Event {
	allowed := m.buildAllowSet()
	var events []Event

	for _, s := range states {
		key := portKey(s.Port, s.Protocol)
		if _, ok := allowed[key]; !ok {
			events = append(events, NewEvent(
				LevelAlert,
				s.Port,
				s.Protocol,
				s.Process,
				"unexpected listener detected",
			))
		}
	}
	return events
}

// IsAllowed reports whether the given port and protocol combination is
// permitted by the current configuration (i.e. has a non-alerting rule).
func (m *Matcher) IsAllowed(port int, protocol string) bool {
	allowed := m.buildAllowSet()
	_, ok := allowed[portKey(port, protocol)]
	return ok
}

func (m *Matcher) buildAllowSet() map[string]struct{} {
	set := make(map[string]struct{})
	for _, r := range m.cfg.Rules {
		if !r.Alert {
			for _, p := range r.Ports {
				set[portKey(p, r.Protocol)] = struct{}{}
			}
		}
	}
	return set
}

func portKey(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
