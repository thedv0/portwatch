package snapshot

import (
	"testing"

	"github.com/netwatch/portwatch/internal/scanner"
)

func pipePort(port int, proto, proc string, pid int) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Process: proc, PID: pid}
}

func TestRunPipeline_AllStagesEnabled(t *testing.T) {
	ports := []scanner.Port{
		pipePort(80, "TCP", "  nginx  ", 100),
		pipePort(80, "TCP", "nginx", 100), // duplicate
	}
	opts := DefaultPipelineOptions()
	result := RunPipeline(ports, opts)

	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port after dedupe, got %d", len(result.Ports))
	}
	if result.Ports[0].Protocol != "tcp" {
		t.Errorf("expected normalized protocol 'tcp', got %q", result.Ports[0].Protocol)
	}
	if result.Ports[0].Process != "nginx" {
		t.Errorf("expected trimmed process 'nginx', got %q", result.Ports[0].Process)
	}
}

func TestRunPipeline_ValidationPopulated(t *testing.T) {
	ports := []scanner.Port{
		pipePort(80, "tcp", "nginx", 1),
	}
	opts := DefaultPipelineOptions()
	result := RunPipeline(ports, opts)

	if result.Validation == nil {
		t.Fatal("expected validation result, got nil")
	}
}

func TestRunPipeline_ClassifiedPopulated(t *testing.T) {
	ports := []scanner.Port{
		pipePort(443, "tcp", "nginx", 1),
	}
	opts := DefaultPipelineOptions()
	result := RunPipeline(ports, opts)

	if len(result.Classified) == 0 {
		t.Fatal("expected classified ports, got none")
	}
}

func TestRunPipeline_NoStages(t *testing.T) {
	ports := []scanner.Port{
		pipePort(80, "TCP", "  svc  ", 10),
		pipePort(80, "TCP", "  svc  ", 10),
	}
	opts := PipelineOptions{} // all disabled
	result := RunPipeline(ports, opts)

	if len(result.Ports) != 2 {
		t.Errorf("expected 2 ports with no stages, got %d", len(result.Ports))
	}
	if result.Validation != nil {
		t.Error("expected nil validation when disabled")
	}
	if len(result.Classified) != 0 {
		t.Error("expected empty classified when disabled")
	}
}

func TestRunPipeline_DoesNotMutateInput(t *testing.T) {
	original := []scanner.Port{
		pipePort(22, "TCP", "sshd", 55),
	}
	originalProtocol := original[0].Protocol
	RunPipeline(original, DefaultPipelineOptions())

	if original[0].Protocol != originalProtocol {
		t.Error("RunPipeline mutated the input slice")
	}
}

func TestDefaultPipelineOptions_AllEnabled(t *testing.T) {
	opts := DefaultPipelineOptions()
	if !opts.Normalize || !opts.Dedupe || !opts.Enrich || !opts.Validate || !opts.Classify {
		t.Error("expected all stages enabled by default")
	}
}
