package snapshot

// FilterOptions controls which ports are included in results.
type FilterOptions struct {
	Protocol  string // empty means all
	MinPort   int    // 0 means no lower bound
	MaxPort   int    // 0 means no upper bound
	PIDZero   bool   // if true, only include entries with PID == 0
	Processes []string // if non-empty, only include matching process names
}

// Filter returns a subset of ports matching the given options.
func Filter(ports []PortState, opts FilterOptions) []PortState {
	result := make([]PortState, 0, len(ports))
	for _, p := range ports {
		if opts.Protocol != "" && p.Protocol != opts.Protocol {
			continue
		}
		if opts.MinPort > 0 && p.Port < opts.MinPort {
			continue
		}
		if opts.MaxPort > 0 && p.Port > opts.MaxPort {
			continue
		}
		if opts.PIDZero && p.PID != 0 {
			continue
		}
		if len(opts.Processes) > 0 && !containsProcess(opts.Processes, p.Process) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func containsProcess(list []string, name string) bool {
	for _, p := range list {
		if p == name {
			return true
		}
	}
	return false
}
