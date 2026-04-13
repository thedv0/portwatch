package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ForecastReport is the serialisable report produced from a forecast result.
type ForecastReport struct {
	GeneratedAt time.Time               `json:"generated_at"`
	Horizon     string                  `json:"horizon"`
	Slope       float64                 `json:"slope"`
	Intercept   float64                 `json:"intercept"`
	Points      []snapshot.ForecastPoint `json:"points"`
}

// BuildForecastReport converts a ForecastResult into a ForecastReport.
func BuildForecastReport(res snapshot.ForecastResult) ForecastReport {
	return ForecastReport{
		GeneratedAt: res.GeneratedAt,
		Horizon:     res.Horizon.String(),
		Slope:       res.Slope,
		Intercept:   res.Intercept,
		Points:      res.Points,
	}
}

// WriteForecastText writes a human-readable forecast report to w.
func WriteForecastText(w io.Writer, r ForecastReport) error {
	fmt.Fprintf(w, "Forecast Report\n")
	fmt.Fprintf(w, "Generated : %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Horizon   : %s\n", r.Horizon)
	fmt.Fprintf(w, "Slope     : %.4f ports/sec\n", r.Slope)
	fmt.Fprintf(w, "Intercept : %.4f\n", r.Intercept)
	fmt.Fprintf(w, "\n%-30s  %s\n", "Time", "Predicted Ports")
	fmt.Fprintf(w, "%s\n", "---------------------------------------------")
	for _, p := range r.Points {
		fmt.Fprintf(w, "%-30s  %.2f\n", p.At.Format(time.RFC3339), p.PortCount)
	}
	return nil
}

// WriteForecastJSON writes the forecast report as JSON to w.
func WriteForecastJSON(w io.Writer, r ForecastReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
