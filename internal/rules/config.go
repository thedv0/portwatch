package rules

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Rule defines a set of ports and the action to take when they are detected.
type Rule struct {
	Name     string `yaml:"name"`
	Ports    []int  `yaml:"ports"`
	Protocol string `yaml:"protocol"`
	Alert    bool   `yaml:"alert"`
}

// Config holds the top-level portwatch configuration.
type Config struct {
	Interval time.Duration `yaml:"interval"`
	LogLevel string        `yaml:"log_level"`
	Rules    []Rule        `yaml:"rules"`
}

// DefaultConfig returns a Config populated with safe defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval: 30 * time.Second,
		LogLevel: "info",
	}
}

// LoadConfig reads and parses a YAML config file at path.
// Missing optional fields fall back to DefaultConfig values.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks Config for logical errors.
func (c *Config) Validate() error {
	if c.Interval < time.Second {
		return fmt.Errorf("interval must be at least 1s, got %s", c.Interval)
	}
	for i, r := range c.Rules {
		if r.Name == "" {
			return fmt.Errorf("rule[%d]: name is required", i)
		}
		if r.Protocol != "tcp" && r.Protocol != "udp" {
			return fmt.Errorf("rule[%d] %q: protocol must be tcp or udp", i, r.Name)
		}
		if len(r.Ports) == 0 {
			return fmt.Errorf("rule[%d] %q: at least one port required", i, r.Name)
		}
	}
	return nil
}
