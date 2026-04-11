package audit

import (
	"fmt"

	"github.com/user/portwatch/internal/snapshot"
)

// DiffAuditor emits audit log entries derived from snapshot diffs.
type DiffAuditor struct {
	logger *Logger
}

// NewDiffAuditor wraps a Logger to audit snapshot diff events.
func NewDiffAuditor(l *Logger) *DiffAuditor {
	return &DiffAuditor{logger: l}
}

// AuditDiff logs added and removed ports from a snapshot.DiffResult.
func (a *DiffAuditor) AuditDiff(diff snapshot.DiffResult) error {
	for _, p := range diff.Added {
		details := map[string]string{
			"port":     fmt.Sprintf("%d", p.Port),
			"protocol": p.Protocol,
			"process":  p.Process,
			"pid":      fmt.Sprintf("%d", p.PID),
		}
		if err := a.logger.Alert("port_opened", details); err != nil {
			return err
		}
	}
	for _, p := range diff.Removed {
		details := map[string]string{
			"port":     fmt.Sprintf("%d", p.Port),
			"protocol": p.Protocol,
			"process":  p.Process,
			"pid":      fmt.Sprintf("%d", p.PID),
		}
		if err := a.logger.Info("port_closed", details); err != nil {
			return err
		}
	}
	return nil
}
