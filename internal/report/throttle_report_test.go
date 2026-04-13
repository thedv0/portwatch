package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeThrottleData() (map[string]time.Time, map[string]int) {
	allowed := map[string]time.Time{
		"tcp:22":  time.Unix(1700000000, 0).UTC(),
		"tcp:443": time.Unix(1700000100, 0).UTC(),
	}
	blocked := map[string]int{
		"tcp:22": 5,
	}
	return allowed, blocked
}

func TestBuildThrottleReport_Counts(t *testing.T) {
	allowed, blocked := makeThrottleData()
	r := BuildThrottleReport(allowed, blocked)
	if r.TotalAllowed != 2 {
		t.Errorf("expected TotalAllowed=2, got %d", r.TotalAllowed)
	}
	if r.TotalBlocked != 5 {
		t.Errorf("expected TotalBlocked=5, got %d", r.TotalBlocked)
	}
	if len(r.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestBuildThrottleReport_TimestampSet(t *testing.T) {
	r := BuildThrottleReport(map[string]time.Time{"udp:53": time.Now()}, nil)
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestBuildThrottleReport_EmptyInput(t *testing.T) {
	r := BuildThrottleReport(nil, nil)
	if r.TotalAllowed != 0 || r.TotalBlocked != 0 {
		t.Error("expected zero counts for empty input")
	}
}

func TestWriteThrottleText_ContainsHeaders(t *testing.T) {
	allowed, blocked := makeThrottleData()
	r := BuildThrottleReport(allowed, blocked)
	var buf bytes.Buffer
	if err := WriteThrottleText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Throttle Report", "Allowed", "Blocked", "tcp:22"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteThrottleText_NoEntries(t *testing.T) {
	r := BuildThrottleReport(nil, nil)
	var buf bytes.Buffer
	WriteThrottleText(&buf, r)
	if !strings.Contains(buf.String(), "no entries") {
		t.Error("expected 'no entries' message for empty report")
	}
}

func TestWriteThrottleJSON_ValidJSON(t *testing.T) {
	allowed, blocked := makeThrottleData()
	r := BuildThrottleReport(allowed, blocked)
	var buf bytes.Buffer
	if err := WriteThrottleJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out ThrottleReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalAllowed != r.TotalAllowed {
		t.Errorf("JSON round-trip: expected TotalAllowed=%d, got %d", r.TotalAllowed, out.TotalAllowed)
	}
}
