package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/netwatch/portwatch/internal/scanner"
)

// ClusterReport holds a rendered cluster result.
type ClusterReport struct {
	Timestamp   time.Time                      `json:"timestamp"`
	GroupBy     string                         `json:"group_by"`
	ClusterCount int                           `json:"cluster_count"`
	TotalPorts  int                            `json:"total_ports"`
	Clusters    map[string][]scanner.Port      `json:"clusters"`
}

// BuildClusterReport constructs a ClusterReport from the raw cluster map.
func BuildClusterReport(clusters map[string][]scanner.Port, groupBy string) ClusterReport {
	total := 0
	for _, v := range clusters {
		total += len(v)
	}
	return ClusterReport{
		Timestamp:    time.Now().UTC(),
		GroupBy:      groupBy,
		ClusterCount: len(clusters),
		TotalPorts:   total,
		Clusters:     clusters,
	}
}

// WriteClusterText writes a human-readable cluster report to w.
func WriteClusterText(w io.Writer, r ClusterReport) error {
	fmt.Fprintf(w, "Cluster Report  [%s]\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Grouped by : %s\n", r.GroupBy)
	fmt.Fprintf(w, "Clusters   : %d\n", r.ClusterCount)
	fmt.Fprintf(w, "Total ports: %d\n\n", r.TotalPorts)

	keys := make([]string, 0, len(r.Clusters))
	for k := range r.Clusters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		ports := r.Clusters[k]
		fmt.Fprintf(w, "  [%s] (%d ports)\n", k, len(ports))
		for _, p := range ports {
			fmt.Fprintf(w, "    %-6d %-5s %s (pid %d)\n", p.Port, p.Protocol, p.Process, p.PID)
		}
	}
	return nil
}

// WriteClusterJSON writes the report as JSON to w.
func WriteClusterJSON(w io.Writer, r ClusterReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
