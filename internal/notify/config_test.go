package notify

import (
	"bytes"
	"testing"
)

func TestBuildDispatcher_DefaultsToLog(t *testing.T) {
	d, err := BuildDispatcher(nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.channels) != 1 {
		t.Errorf("expected 1 default channel, got %d", len(d.channels))
	}
}

func TestBuildDispatcher_LogType(t *testing.T) {
	cfgs := []ChannelConfig{{Type: "log"}}
	var buf bytes.Buffer
	d, err := BuildDispatcher(cfgs, &buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil dispatcher")
	}
}

func TestBuildDispatcher_ExecMissingCommand(t *testing.T) {
	cfgs := []ChannelConfig{{Type: "exec"}}
	_, err := BuildDispatcher(cfgs, nil, nil)
	if err == nil {
		t.Fatal("expected error for exec without command")
	}
}

func TestBuildDispatcher_UnknownType(t *testing.T) {
	cfgs := []ChannelConfig{{Type: "webhook"}}
	_, err := BuildDispatcher(cfgs, nil, nil)
	if err == nil {
		t.Fatal("expected error for unknown channel type")
	}
}

func TestBuildDispatcher_MultipleChannels(t *testing.T) {
	cfgs := []ChannelConfig{
		{Type: "log"},
		{Type: "exec", Command: "echo", Args: []string{"hi"}},
	}
	var buf bytes.Buffer
	d, err := BuildDispatcher(cfgs, &buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(d.channels))
	}
}
