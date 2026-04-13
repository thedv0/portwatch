package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/portwatch/internal/snapshot"
)

// AnomalyReport holds a summary of detected anomalies.
type AnomalyReport struct {
	Timestamp  time.Time          `json:"timestamp"`
	Total      int                `json:"total"`
	NewPorts   int                `json:"new_ports"`
	GonePorts  int                `json:"gone_ports"`
	PIDChanges int                `json:"pid_changes"`
	Spikes     int                `json:"spikes"`
	Anomalies  []snapshot.Anomaly `json:"anomalies"`
}

// BuildAnomalyReport constructs an AnomalyReport from a list of anomalies.
func BuildAnomalyReport(anomalies []snapshot.Anomaly, now time.Time) AnomalyReport {
	r := AnomalyReport{
		Timestamp: now,
		Total:     len(anomalies),
		Anomalies: anomalies,
	}
	for _, a := range anomalies {
		switch a.Type {
		case snapshot.AnomalyNewPort:
			r.NewPorts++
		case snapshot.AnomalyGonePort:
			r.GonePorts++
		case snapshot.AnomalyPIDChanged:
			r.PIDChanges++
		case snapshot.AnomalyPortSpike:
			r.Spikes++
		}
	}
	return r
}

// WriteAnomalyText writes a human-readable anomaly report to w.
func WriteAnomalyText(w io.Writer, r AnomalyReport) error {
	fmt.Fprintf(w, "Anomaly Report — %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Total: %d  New: %d  Gone: %d  PID Changes: %d  Spikes: %d\n",
		r.Total, r.NewPorts, r.GonePorts, r.PIDChanges, r.Spikes)
	if len(r.Anomalies) == 0 {
		fmt.Fprintln(w, "No anomalies detected.")
		return nil
	}
	fmt.Fprintln(w, "---")
	for _, a := range r.Anomalies {
		fmt.Fprintf(w, "[%s] port=%d proto=%s process=%s pid=%d — %s\n",
			a.Type, a.Port, a.Protocol, a.Process, a.PID, a.Message)
	}
	return nil
}

// WriteAnomalyJSON writes the report as JSON to w.
func WriteAnomalyJSON(w io.Writer, r AnomalyReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
