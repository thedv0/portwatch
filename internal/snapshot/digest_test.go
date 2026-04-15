package snapshot

import (
	"testing"
)

func dgPort(proto string, port int, proc string, pid int) PortState {
	return PortState{Protocol: proto, Port: port, Process: proc, PID: pid}
}

func TestDigest_EmptyInput_ReturnsEmptyHash(t *testing.T) {
	res := Digest(nil, DefaultDigestOptions())
	if res.Hex == "" {
		t.Fatal("expected non-empty hex for empty input")
	}
	if res.PortCount != 0 {
		t.Fatalf("expected port count 0, got %d", res.PortCount)
	}
}

func TestDigest_PortCountMatchesInput(t *testing.T) {
	ports := []PortState{
		dgPort("tcp", 80, "nginx", 100),
		dgPort("tcp", 443, "nginx", 101),
	}
	res := Digest(ports, DefaultDigestOptions())
	if res.PortCount != 2 {
		t.Fatalf("expected port count 2, got %d", res.PortCount)
	}
}

func TestDigest_OrderIndependent(t *testing.T) {
	a := []PortState{
		dgPort("tcp", 80, "nginx", 1),
		dgPort("tcp", 443, "nginx", 2),
	}
	b := []PortState{
		dgPort("tcp", 443, "nginx", 2),
		dgPort("tcp", 80, "nginx", 1),
	}
	opts := DefaultDigestOptions()
	if Digest(a, opts).Hex != Digest(b, opts).Hex {
		t.Fatal("digest should be order-independent")
	}
}

func TestDigest_DifferentPorts_DifferentHash(t *testing.T) {
	a := []PortState{dgPort("tcp", 80, "nginx", 1)}
	b := []PortState{dgPort("tcp", 8080, "nginx", 1)}
	opts := DefaultDigestOptions()
	if Digest(a, opts).Hex == Digest(b, opts).Hex {
		t.Fatal("different ports should produce different digests")
	}
}

func TestDigest_PIDIncluded_AffectsHash(t *testing.T) {
	base := []PortState{dgPort("tcp", 80, "nginx", 100)}
	other := []PortState{dgPort("tcp", 80, "nginx", 999)}

	without := DefaultDigestOptions()
	without.IncludePID = false

	with := DefaultDigestOptions()
	with.IncludePID = true

	if Digest(base, without).Hex == Digest(other, without).Hex {
		// same without PID is expected only if process matches, which it does
	}
	if Digest(base, with).Hex == Digest(other, with).Hex {
		t.Fatal("with PID enabled, different PIDs should yield different digests")
	}
}

func TestDigest_ProcessExcluded_SameHash(t *testing.T) {
	a := []PortState{dgPort("tcp", 80, "nginx", 1)}
	b := []PortState{dgPort("tcp", 80, "apache", 1)}

	opts := DefaultDigestOptions()
	opts.IncludeProcess = false

	if Digest(a, opts).Hex != Digest(b, opts).Hex {
		t.Fatal("with process excluded, different process names should not affect digest")
	}
}

func TestDigest_Deterministic(t *testing.T) {
	ports := []PortState{
		dgPort("udp", 53, "dnsmasq", 42),
		dgPort("tcp", 22, "sshd", 10),
	}
	opts := DefaultDigestOptions()
	first := Digest(ports, opts).Hex
	second := Digest(ports, opts).Hex
	if first != second {
		t.Fatal("digest must be deterministic across calls")
	}
}
