package snapshot

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeClassPort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, PID: 1, Process: "test"}
}

func TestClassify_SystemPort(t *testing.T) {
	ports := []scanner.Port{makeClassPort(80, "tcp")}
	result := Classify(ports, DefaultClassifyOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Class != RiskClassSystem {
		t.Errorf("expected system class, got %s", result[0].Class)
	}
}

func TestClassify_RegisteredPort(t *testing.T) {
	ports := []scanner.Port{makeClassPort(8080, "tcp")}
	result := Classify(ports, DefaultClassifyOptions())
	if result[0].Class != RiskClassRegistered {
		t.Errorf("expected registered class, got %s", result[0].Class)
	}
}

func TestClassify_DynamicPort(t *testing.T) {
	ports := []scanner.Port{makeClassPort(55000, "tcp")}
	result := Classify(ports, DefaultClassifyOptions())
	if result[0].Class != RiskClassDynamic {
		t.Errorf("expected dynamic class, got %s", result[0].Class)
	}
}

func TestClassify_WellKnownFlag(t *testing.T) {
	ports := []scanner.Port{
		makeClassPort(443, "tcp"),
		makeClassPort(9999, "tcp"),
	}
	result := Classify(ports, DefaultClassifyOptions())
	if !result[0].WellKnown {
		t.Error("expected port 443 to be well-known")
	}
	if result[1].WellKnown {
		t.Error("expected port 9999 to not be well-known")
	}
}

func TestClassify_EmptyInput(t *testing.T) {
	result := Classify(nil, DefaultClassifyOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestClassify_CustomWellKnown(t *testing.T) {
	opts := ClassifyOptions{WellKnownPorts: []int{9200}}
	ports := []scanner.Port{makeClassPort(9200, "tcp")}
	result := Classify(ports, opts)
	if !result[0].WellKnown {
		t.Error("expected port 9200 to be well-known with custom opts")
	}
}

func TestDefaultClassifyOptions_HasDefaults(t *testing.T) {
	opts := DefaultClassifyOptions()
	if len(opts.WellKnownPorts) == 0 {
		t.Error("expected non-empty default well-known ports")
	}
}
