package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	r := New()
	r.Counter("scans").Add(3)
	r.Gauge("open_ports").Set(12)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	Handler(r).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}
	var snap map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestHandler_ContainsCounterAndGauge(t *testing.T) {
	r := New()
	r.Counter("alerts_sent").Add(5)
	r.Gauge("open_ports").Set(8)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	Handler(r).ServeHTTP(rec, req)

	var snap map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&snap)

	if snap["alerts_sent"].(float64) != 5 {
		t.Fatalf("unexpected alerts_sent: %v", snap["alerts_sent"])
	}
	if snap["open_ports"].(float64) != 8 {
		t.Fatalf("unexpected open_ports: %v", snap["open_ports"])
	}
}

func TestHandler_UptimePresent(t *testing.T) {
	r := New()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	Handler(r).ServeHTTP(rec, req)

	var snap map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&snap)
	if _, ok := snap["uptime_seconds"]; !ok {
		t.Fatal("uptime_seconds missing from response")
	}
}
