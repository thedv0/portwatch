package snapshot

import "time"

// AnomalyType categorizes the kind of anomaly detected.
type AnomalyType string

const (
	AnomalyNewPort    AnomalyType = "new_port"
	AnomalyGonePort   AnomalyType = "gone_port"
	AnomalyPortSpike  AnomalyType = "port_spike"
	AnomalyPIDChanged AnomalyType = "pid_changed"
)

// Anomaly represents a detected anomaly in port state.
type Anomaly struct {
	Type      AnomalyType `json:"type"`
	Port      int         `json:"port"`
	Protocol  string      `json:"protocol"`
	Process   string      `json:"process,omitempty"`
	PID       int         `json:"pid,omitempty"`
	Message   string      `json:"message"`
	DetectedAt time.Time  `json:"detected_at"`
}

// DefaultAnomalyOptions returns sensible defaults.
type AnomalyOptions struct {
	SpikeThreshold int // number of new ports in one cycle to trigger spike
	TrackPIDChange bool
}

func DefaultAnomalyOptions() AnomalyOptions {
	return AnomalyOptions{
		SpikeThreshold: 5,
		TrackPIDChange: true,
	}
}

// DetectAnomalies compares previous and current port states and returns anomalies.
func DetectAnomalies(prev, curr []PortState, opts AnomalyOptions, now time.Time) []Anomaly {
	var anomalies []Anomaly

	prevIdx := indexByKey(prev)
	currIdx := indexByKey(curr)

	newCount := 0
	for k, cp := range currIdx {
		if _, exists := prevIdx[k]; !exists {
			newCount++
			anomalies = append(anomalies, Anomaly{
				Type:       AnomalyNewPort,
				Port:       cp.Port,
				Protocol:   cp.Protocol,
				Process:    cp.Process,
				PID:        cp.PID,
				Message:    "new port opened",
				DetectedAt: now,
			})
		} else if opts.TrackPIDChange {
			pp := prevIdx[k]
			if pp.PID != cp.PID && cp.PID != 0 && pp.PID != 0 {
				anomalies = append(anomalies, Anomaly{
					Type:       AnomalyPIDChanged,
					Port:       cp.Port,
					Protocol:   cp.Protocol,
					Process:    cp.Process,
					PID:        cp.PID,
					Message:    "PID changed on port",
					DetectedAt: now,
				})
			}
		}
	}

	for k, pp := range prevIdx {
		if _, exists := currIdx[k]; !exists {
			anomalies = append(anomalies, Anomaly{
				Type:       AnomalyGonePort,
				Port:       pp.Port,
				Protocol:   pp.Protocol,
				Process:    pp.Process,
				PID:        pp.PID,
				Message:    "port closed",
				DetectedAt: now,
			})
		}
	}

	if newCount >= opts.SpikeThreshold {
		anomalies = append(anomalies, Anomaly{
			Type:       AnomalyPortSpike,
			Message:    "sudden spike in new open ports",
			DetectedAt: now,
		})
	}

	return anomalies
}

func indexByKey(ports []PortState) map[string]PortState {
	m := make(map[string]PortState, len(ports))
	for _, p := range ports {
		m[anomalyKey(p)] = p
	}
	return m
}

func anomalyKey(p PortState) string {
	return p.Protocol + ":" + itoa(p.Port)
}
