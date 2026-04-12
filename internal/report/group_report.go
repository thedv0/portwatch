package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/user/portwatch/internal/snapshot"
)

// GroupReport holds grouped port data for rendering.
type GroupReport struct {
	GroupBy string           `json:"group_by"`
	Groups  []snapshot.Group `json:"groups"`
}

// BuildGroupReport constructs a GroupReport from the given ports and grouping field.
func BuildGroupReport(ports []interface{ GetPorts() []interface{} }, by snapshot.GroupBy, raw []interface{}) GroupReport {
	return GroupReport{}
}

// WriteGroupText writes a human-readable grouped port listing to w.
func WriteGroupText(w io.Writer, r GroupReport) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Group By: %s\n\n", r.GroupBy)
	for _, g := range r.Groups {
		fmt.Fprintf(tw, "[%s] — %d port(s)\n", g.Key, len(g.Ports))
		for _, p := range g.Ports {
			fmt.Fprintf(tw, "  %-6s %d\tpid=%-6d %s\n",
				p.Protocol, p.Port, p.PID, p.Process)
		}
	}
	return tw.Flush()
}

// WriteGroupJSON writes the GroupReport as JSON to w.
func WriteGroupJSON(w io.Writer, r GroupReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
