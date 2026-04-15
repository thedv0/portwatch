package snapshot

import (
	"testing"
	"time"
)

func sealPort(port int, proto, process string, pid int) Port {
	return Port{Port: port, Protocol: proto, Process: process, PID: pid}
}

func fixedSealClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestSeal_ProducesNonEmptyHash(t *testing.T) {
	snap := Snapshot{Ports: []Port{sealPort(80, "tcp", "nginx", 100)}}
	sealed, err := Seal(snap, DefaultSealOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sealed.Hash == "" {
		t.Error("expected non-empty hash")
	}
}

func TestSeal_SealedAtSet(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	opts := DefaultSealOptions()
	opts.Clock = fixedSealClock(now)
	snap := Snapshot{Ports: []Port{sealPort(443, "tcp", "nginx", 200)}}
	sealed, err := Seal(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sealed.SealedAt.Equal(now) {
		t.Errorf("expected SealedAt %v, got %v", now, sealed.SealedAt)
	}
}

func TestVerify_ValidSeal_NoError(t *testing.T) {
	snap := Snapshot{Ports: []Port{sealPort(22, "tcp", "sshd", 50)}}
	opts := DefaultSealOptions()
	sealed, err := Seal(snap, opts)
	if err != nil {
		t.Fatalf("seal error: %v", err)
	}
	if err := Verify(sealed, opts); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestVerify_TamperedPorts_ReturnsError(t *testing.T) {
	snap := Snapshot{Ports: []Port{sealPort(22, "tcp", "sshd", 50)}}
	opts := DefaultSealOptions()
	sealed, err := Seal(snap, opts)
	if err != nil {
		t.Fatalf("seal error: %v", err)
	}
	// Tamper with the snapshot after sealing.
	sealed.Snapshot.Ports = append(sealed.Snapshot.Ports, sealPort(8080, "tcp", "evil", 9999))
	if err := Verify(sealed, opts); err == nil {
		t.Error("expected error for tampered snapshot, got nil")
	}
}

func TestSeal_DeterministicHash(t *testing.T) {
	snap := Snapshot{Ports: []Port{sealPort(80, "tcp", "nginx", 1)}}
	opts := DefaultSealOptions()
	s1, _ := Seal(snap, opts)
	s2, _ := Seal(snap, opts)
	if s1.Hash != s2.Hash {
		t.Errorf("expected deterministic hash, got %q and %q", s1.Hash, s2.Hash)
	}
}

func TestSeal_IncludeTimestamp_DifferentHashes(t *testing.T) {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	ports := []Port{sealPort(80, "tcp", "nginx", 1)}
	opts := SealOptions{IncludeTimestamp: true, Clock: fixedSealClock(t1)}
	s1, _ := Seal(Snapshot{Ports: ports, Timestamp: t1}, opts)
	s2, _ := Seal(Snapshot{Ports: ports, Timestamp: t2}, opts)
	if s1.Hash == s2.Hash {
		t.Error("expected different hashes for different timestamps when IncludeTimestamp=true")
	}
}
