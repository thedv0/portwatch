package snapshot

import (
	"testing"
	"time"
)

func expPort(port int, lastSeen time.Time) PortState {
	return PortState{Port: port, Protocol: "tcp", LastSeen: lastSeen}
}

func TestExpire_RetainsRecentPorts(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	opts := DefaultExpireOptions()
	opts.MaxAge = 5 * time.Minute
	opts.Now = func() time.Time { return now }

	ports := []PortState{
		expPort(80, now.Add(-1*time.Minute)),
		expPort(443, now.Add(-3*time.Minute)),
	}
	res := Expire(ports, opts)
	if len(res.Retained) != 2 {
		t.Fatalf("expected 2 retained, got %d", len(res.Retained))
	}
	if len(res.Expired) != 0 {
		t.Fatalf("expected 0 expired, got %d", len(res.Expired))
	}
}

func TestExpire_RemovesOldPorts(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	opts := DefaultExpireOptions()
	opts.MaxAge = 5 * time.Minute
	opts.Now = func() time.Time { return now }

	ports := []PortState{
		expPort(22, now.Add(-10*time.Minute)),
		expPort(80, now.Add(-1*time.Minute)),
	}
	res := Expire(ports, opts)
	if len(res.Retained) != 1 {
		t.Fatalf("expected 1 retained, got %d", len(res.Retained))
	}
	if res.Retained[0].Port != 80 {
		t.Errorf("expected port 80 retained, got %d", res.Retained[0].Port)
	}
	if len(res.Expired) != 1 || res.Expired[0].Port != 22 {
		t.Errorf("expected port 22 expired")
	}
}

func TestExpire_KeepUnseenPorts_True(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	opts := DefaultExpireOptions()
	opts.Now = func() time.Time { return now }
	opts.KeepUnseenPorts = true

	ports := []PortState{expPort(9000, time.Time{})}
	res := Expire(ports, opts)
	if len(res.Retained) != 1 {
		t.Errorf("expected unseen port retained")
	}
}

func TestExpire_KeepUnseenPorts_False(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	opts := DefaultExpireOptions()
	opts.Now = func() time.Time { return now }
	opts.KeepUnseenPorts = false

	ports := []PortState{expPort(9000, time.Time{})}
	res := Expire(ports, opts)
	if len(res.Expired) != 1 {
		t.Errorf("expected unseen port expired")
	}
}

func TestExpire_TotalCount(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	opts := DefaultExpireOptions()
	opts.MaxAge = 5 * time.Minute
	opts.Now = func() time.Time { return now }

	ports := []PortState{
		expPort(80, now.Add(-1*time.Minute)),
		expPort(22, now.Add(-10*time.Minute)),
		expPort(443, now.Add(-2*time.Minute)),
	}
	res := Expire(ports, opts)
	if res.Total != 3 {
		t.Errorf("expected total 3, got %d", res.Total)
	}
}

func TestExpire_EmptyInput(t *testing.T) {
	opts := DefaultExpireOptions()
	res := Expire(nil, opts)
	if res.Total != 0 || len(res.Retained) != 0 || len(res.Expired) != 0 {
		t.Errorf("expected empty result for nil input")
	}
}
