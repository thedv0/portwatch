package notify

import (
	"fmt"
	"io"
	"log"
	"os"
)

// ChannelConfig describes a single notification channel from YAML config.
type ChannelConfig struct {
	Type    string   `yaml:"type"`    // "log" or "exec"
	Command string   `yaml:"command"` // for exec type
	Args    []string `yaml:"args"`    // for exec type
}

// BuildDispatcher constructs a Dispatcher from a slice of ChannelConfig.
// If no channels are configured, a default log channel writing to logOut is used.
// If logOut is nil, os.Stderr is used.
func BuildDispatcher(cfgs []ChannelConfig, logOut io.Writer, logger *log.Logger) (*Dispatcher, error) {
	if logOut == nil {
		logOut = os.Stderr
	}

	var channels []Channel
	for i, cfg := range cfgs {
		switch cfg.Type {
		case "log", "":
			channels = append(channels, NewLogChannel(logOut))
		case "exec":
			if cfg.Command == "" {
				return nil, fmt.Errorf("channel[%d]: exec channel requires a command", i)
			}
			channels = append(channels, NewExecChannel(cfg.Command, cfg.Args))
		default:
			return nil, fmt.Errorf("channel[%d]: unknown channel type %q", i, cfg.Type)
		}
	}

	// Always ensure at least one channel.
	if len(channels) == 0 {
		channels = append(channels, NewLogChannel(logOut))
	}

	return NewDispatcher(logger, channels...), nil
}
