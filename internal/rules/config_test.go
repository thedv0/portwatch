package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	content := `
version: "1"
interval_seconds: 10
rules:
  - name: allow-ssh
    ports: ["22"]
    protocols: ["tcp"]
    action: allow
  - name: alert-unknown
    ports: ["8080", "9000-"]
    action: alert
`
	path := writeTemp(t, content)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 10 {
		t.Errorf("expected interval 10, got %d", cfg.Interval)
	}
	if len(cfg.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cfg.Rules))
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 30 {
		t.Errorf("expected default interval 30, got %d", cfg.Interval)
	}
	if cfg.Version != "1" {
		t.Errorf("expected default version '1', got %s", cfg.Version)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_InvalidInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for interval=0")
	}
}

func TestValidate_MissingRuleName(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []Rule{{Ports: []string{"80"}, Action: ActionAlert}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing rule name")
	}
}

func TestValidate_UnknownAction(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []Rule{{Name: "test", Ports: []string{"80"}, Action: "deny"}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestValidate_EmptyPorts(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []Rule{{Name: "test", Ports: []string{}, Action: ActionAllow}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty ports")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []Rule{
		{Name: "allow-ssh", Ports: []string{"22"}, Protocols: []string{"tcp"}, Action: ActionAllow},
		{Name: "alert-http", Ports: []string{"80", "8080"}, Action: ActionAlert},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no validation error, got: %v", err)
	}
}
