package audit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileConfig configures a file-backed audit log.
type FileConfig struct {
	Path    string `yaml:"path"`
	Enabled bool   `yaml:"enabled"`
}

// OpenFile opens or creates the audit log file at cfg.Path.
// The caller is responsible for closing the returned WriteCloser.
func OpenFile(cfg FileConfig) (io.WriteCloser, error) {
	if !cfg.Enabled || cfg.Path == "" {
		return nopCloser{io.Discard}, nil
	}
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir %s: %w", dir, err)
	}
	f, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file %s: %w", cfg.Path, err)
	}
	return f, nil
}

// NewFileLogger creates a Logger backed by the file at cfg.Path.
// Returns a Logger and a cleanup function to close the underlying file.
func NewFileLogger(cfg FileConfig) (*Logger, func() error, error) {
	wc, err := OpenFile(cfg)
	if err != nil {
		return nil, nil, err
	}
	return NewLogger(wc), wc.Close, nil
}

// nopCloser wraps an io.Writer with a no-op Close.
type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }
