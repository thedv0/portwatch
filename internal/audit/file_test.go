package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenFile_Disabled_ReturnsDiscard(t *testing.T) {
	wc, err := OpenFile(FileConfig{Enabled: false, Path: "/tmp/should-not-exist.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer wc.Close()
	// writing to discard should not error
	_, err = wc.Write([]byte("hello"))
	if err != nil {
		t.Errorf("write to discard failed: %v", err)
	}
}

func TestOpenFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "audit.log")
	wc, err := OpenFile(FileConfig{Enabled: true, Path: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer wc.Close()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to be created at %s", path)
	}
}

func TestNewFileLogger_WritesToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	l, close, err := NewFileLogger(FileConfig{Enabled: true, Path: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer close()
	if err := l.Alert("test_event", map[string]string{"k": "v"}); err != nil {
		t.Fatalf("log error: %v", err)
	}
	_ = close()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(data), "test_event") {
		t.Errorf("expected test_event in log, got: %s", string(data))
	}
}

func TestOpenFile_EmptyPath_ReturnsDiscard(t *testing.T) {
	wc, err := OpenFile(FileConfig{Enabled: true, Path: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer wc.Close()
	_, err = wc.Write([]byte("data"))
	if err != nil {
		t.Errorf("unexpected write error: %v", err)
	}
}
