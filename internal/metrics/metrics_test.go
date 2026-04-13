package metrics

import (
	"testing"
	"time"
)

func TestCounter_IncAndValue(t *testing.T) {
	c := &Counter{}
	c.Inc()
	c.Inc()
	c.Add(3)
	if got := c.Value(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestCounter_Reset(t *testing.T) {
	c := &Counter{}
	c.Add(10)
	c.Reset()
	if c.Value() != 0 {
		t.Fatal("expected 0 after reset")
	}
}

func TestGauge_SetAndValue(t *testing.T) {
	g := &Gauge{}
	g.Set(3.14)
	if g.Value() != 3.14 {
		t.Fatalf("expected 3.14, got %f", g.Value())
	}
}

func TestRegistry_CounterIdempotent(t *testing.T) {
	r := New()
	a := r.Counter("scans")
	b := r.Counter("scans")
	if a != b {
		t.Fatal("expected same counter instance")
	}
}

func TestRegistry_GaugeIdempotent(t *testing.T) {
	r := New()
	a := r.Gauge("open_ports")
	b := r.Gauge("open_ports")
	if a != b {
		t.Fatal("expected same gauge instance")
	}
}

func TestRegistry_Snapshot_ContainsUptime(t *testing.T) {
	r := New()
	time.Sleep(2 * time.Millisecond)
	snap := r.Snapshot()
	v, ok := snap["uptime_seconds"]
	if !ok {
		t.Fatal("missing uptime_seconds")
	}
	if v.(float64) <= 0 {
		t.Fatal("uptime should be positive")
	}
}

func TestRegistry_Snapshot_ContainsMetrics(t *testing.T) {
	r := New()
	r.Counter("alerts_sent").Add(7)
	r.Gauge("open_ports").Set(42)
	snap := r.Snapshot()
	if snap["alerts_sent"].(int64) != 7 {
		t.Fatalf("unexpected alerts_sent: %v", snap["alerts_sent"])
	}
	if snap["open_ports"].(float64) != 42 {
		t.Fatalf("unexpected open_ports: %v", snap["open_ports"])
	}
}

func TestRegistry_Snapshot_IsolatedBetweenRegistries(t *testing.T) {
	r1 := New()
	r2 := New()
	r1.Counter("hits").Add(10)
	r2.Counter("hits").Add(99)

	snap1 := r1.Snapshot()
	snap2 := r2.Snapshot()

	if snap1["hits"].(int64) != 10 {
		t.Fatalf("r1: expected hits=10, got %v", snap1["hits"])
	}
	if snap2["hits"].(int64) != 99 {
		t.Fatalf("r2: expected hits=99, got %v", snap2["hits"])
	}
}
