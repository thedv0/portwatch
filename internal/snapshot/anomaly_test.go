package snapshot

import (
	"testing"
	"time"
)

func aport(port int, proto, proc string, pid int) PortState {
	return PortState{Port: port, Protocol: proto, Process: proc, PID: pid}
}

var fixedAnomalyTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestDetectAnomalies_NoChange(t *testing.T) {
	ports := []PortState{aport(80, "tcp", "nginx", 100)}
	result := DetectAnomalies(ports, ports, DefaultAnomalyOptions(), fixedAnomalyTime)
	if len(result) != 0 {
		t.Fatalf("expected 0 anomalies, got %d", len(result))
	}
}

func TestDetectAnomalies_NewPort(t *testing.T) {
	prev := []PortState{aport(80, "tcp", "nginx", 100)}
	curr := []PortState{aport(80, "tcp", "nginx", 100), aport(443, "tcp", "nginx", 100)}
	result := DetectAnomalies(prev, curr, DefaultAnomalyOptions(), fixedAnomalyTime)
	if len(result) != 1 || result[0].Type != AnomalyNewPort {
		t.Fatalf("expected 1 new_port anomaly, got %+v", result)
	}
	if result[0].Port != 443 {
		t.Errorf("expected port 443, got %d", result[0].Port)
	}
}

func TestDetectAnomalies_GonePort(t *testing.T) {
	prev := []PortState{aport(80, "tcp", "nginx", 100), aport(8080, "tcp", "app", 200)}
	curr := []PortState{aport(80, "tcp", "nginx", 100)}
	result := DetectAnomalies(prev, curr, DefaultAnomalyOptions(), fixedAnomalyTime)
	if len(result) != 1 || result[0].Type != AnomalyGonePort {
		t.Fatalf("expected 1 gone_port anomaly, got %+v", result)
	}
}

func TestDetectAnomalies_PIDChanged(t *testing.T) {
	prev := []PortState{aport(80, "tcp", "nginx", 100)}
	curr := []PortState{aport(80, "tcp", "nginx", 999)}
	opts := DefaultAnomalyOptions()
	result := DetectAnomalies(prev, curr, opts, fixedAnomalyTime)
	if len(result) != 1 || result[0].Type != AnomalyPIDChanged {
		t.Fatalf("expected 1 pid_changed anomaly, got %+v", result)
	}
}

func TestDetectAnomalies_PIDChangeDisabled(t *testing.T) {
	prev := []PortState{aport(80, "tcp", "nginx", 100)}
	curr := []PortState{aport(80, "tcp", "nginx", 999)}
	opts := DefaultAnomalyOptions()
	opts.TrackPIDChange = false
	result := DetectAnomalies(prev, curr, opts, fixedAnomalyTime)
	if len(result) != 0 {
		t.Errorf("expected 0 anomalies when PID tracking disabled, got %d", len(result))
	}
}

func TestDetectAnomalies_Spike(t *testing.T) {
	prev := []PortState{}
	curr := make([]PortState, 6)
	for i := range curr {
		curr[i] = aport(8000+i, "tcp", "app", i+1)
	}
	opts := DefaultAnomalyOptions()
	opts.SpikeThreshold = 5
	result := DetectAnomalies(prev, curr, opts, fixedAnomalyTime)
	hasSpike := false
	for _, a := range result {
		if a.Type == AnomalyPortSpike {
			hasSpike = true
		}
	}
	if !hasSpike {
		t.Error("expected port_spike anomaly")
	}
}

func TestDetectAnomalies_DetectedAtSet(t *testing.T) {
	prev := []PortState{}
	curr := []PortState{aport(9000, "tcp", "svc", 42)}
	result := DetectAnomalies(prev, curr, DefaultAnomalyOptions(), fixedAnomalyTime)
	if len(result) == 0 {
		t.Fatal("expected anomaly")
	}
	if !result[0].DetectedAt.Equal(fixedAnomalyTime) {
		t.Errorf("expected DetectedAt %v, got %v", fixedAnomalyTime, result[0].DetectedAt)
	}
}
