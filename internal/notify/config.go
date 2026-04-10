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
func BuildDispatcher(cfgs []ChannelConfig, logOut io.Writer, logger *log.Logger) (*Dispatcher, error) {
	if logOut == nil {
		logOut = os.Stderr
	}

	var channels []Channel
	for _, cfg := range cfgs {
		switch cfg.Type {
		case "log", "":
			channels = append(channels, NewLogChannel(logOut))
		case "exec":
			if cfg.Command == "" {
				return nil, fmt.Errorf("exec channel requires a command")
			}
			channels = append(channels, NewExecChannel(cfg.Command, cfg.Args))
		default:
			return nil, fmt.Errorf("unknown channel type %q", cfg.Type)
		}
	}

	// Always ensure at least one channel.
	if len(channels) == 0 {
		channels = append(channels, NewLogChannel(logOut))
	}

	return NewDispatcher(logger, channels...), nil
}
