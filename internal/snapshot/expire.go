package snapshot

import (
	"time"
)

// ExpireOptions controls how port entries are expired from a snapshot.
type ExpireOptions struct {
	// MaxAge is the maximum duration a port entry is considered valid.
	// Entries with a LastSeen time older than MaxAge are expired.
	MaxAge time.Duration

	// Now overrides the current time for testing.
	Now func() time.Time

	// KeepUnseenPorts retains entries that have no LastSeen timestamp (zero time).
	KeepUnseenPorts bool
}

// DefaultExpireOptions returns sensible defaults for ExpireOptions.
func DefaultExpireOptions() ExpireOptions {
	return ExpireOptions{
		MaxAge:          5 * time.Minute,
		Now:             time.Now,
		KeepUnseenPorts: true,
	}
}

// ExpireResult holds the outcome of an Expire operation.
type ExpireResult struct {
	Retained []PortState
	Expired  []PortState
	Total    int
}

// Expire removes port entries from ports that are older than opts.MaxAge.
// Entries with a zero LastSeen are handled according to opts.KeepUnseenPorts.
func Expire(ports []PortState, opts ExpireOptions) ExpireResult {
	if opts.Now == nil {
		opts.Now = time.Now
	}
	now := opts.Now()
	cutoff := now.Add(-opts.MaxAge)

	result := ExpireResult{Total: len(ports)}
	for _, p := range ports {
		if p.LastSeen.IsZero() {
			if opts.KeepUnseenPorts {
				result.Retained = append(result.Retained, p)
			} else {
				result.Expired = append(result.Expired, p)
			}
			continue
		}
		if p.LastSeen.Before(cutoff) {
			result.Expired = append(result.Expired, p)
		} else {
			result.Retained = append(result.Retained, p)
		}
	}
	return result
}
