package rules

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAllow Action = "allow"
	ActionAlert Action = "alert"
	ActionBlock Action = "block"
)

// Rule defines a single port monitoring rule.
type Rule struct {
	Name      string   `yaml:"name"`
	Ports     []string `yaml:"ports"`
	Protocols []string `yaml:"protocols"`
	Action    Action   `yaml:"action"`
	Comment   string   `yaml:"comment,omitempty"`
}

// Config holds the full portwatch configuration.
type Config struct {
	Version  string `yaml:"version"`
	Interval int    `yaml:"interval_seconds"`
	Rules    []Rule `yaml:"rules"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Version:  "1",
		Interval: 30,
		Rules:    []Rule{},
	}
}

// LoadConfig reads and parses a YAML config file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks that the config is semantically valid.
func (c *Config) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval_seconds must be greater than 0")
	}

	for i, rule := range c.Rules {
		if rule.Name == "" {
			return fmt.Errorf("rule[%d]: name is required", i)
		}
		if len(rule.Ports) == 0 {
			return fmt.Errorf("rule[%d] (%s): at least one port is required", i, rule.Name)
		}
		switch rule.Action {
		case ActionAllow, ActionAlert, ActionBlock:
			// valid
		default:
			return fmt.Errorf("rule[%d] (%s): unknown action %q", i, rule.Name, rule.Action)
		}
	}

	return nil
}
