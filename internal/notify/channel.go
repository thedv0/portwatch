package notify

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Channel defines how an alert is dispatched.
type Channel interface {
	Send(e alert.Event) error
	Name() string
}

// LogChannel writes alerts to a writer (default stderr).
type LogChannel struct {
	out  io.Writer
	label string
}

func NewLogChannel(out io.Writer) *LogChannel {
	if out == nil {
		out = os.Stderr
	}
	return &LogChannel{out: out, label: "log"}
}

func (l *LogChannel) Name() string { return l.label }

func (l *LogChannel) Send(e alert.Event) error {
	_, err := fmt.Fprintf(l.out, "[%s] ALERT port=%d proto=%s msg=%s\n",
		e.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		e.Port, e.Protocol, e.Message)
	return err
}

// ExecChannel runs an external command, passing alert details as env vars.
type ExecChannel struct {
	command string
	args    []string
}

func NewExecChannel(command string, args []string) *ExecChannel {
	return &ExecChannel{command: command, args: args}
}

func (e *ExecChannel) Name() string { return "exec:" + e.command }

func (ec *ExecChannel) Send(ev alert.Event) error {
	cmd := exec.Command(ec.command, ec.args...)
	cmd.Env = append(os.Environ(),
		"PORTWATCH_PORT="+fmt.Sprintf("%d", ev.Port),
		"PORTWATCH_PROTO="+ev.Protocol,
		"PORTWATCH_MSG="+ev.Message,
		"PORTWATCH_TS="+ev.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec channel %q failed: %w (output: %s)",
			ec.command, err, strings.TrimSpace(string(out)))
	}
	return nil
}
