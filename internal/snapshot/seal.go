package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SealedSnapshot wraps a snapshot with a content hash for integrity verification.
type SealedSnapshot struct {
	Snapshot  Snapshot  `json:"snapshot"`
	Hash      string    `json:"hash"`
	SealedAt  time.Time `json:"sealed_at"`
}

// DefaultSealOptions returns default options for sealing.
func DefaultSealOptions() SealOptions {
	return SealOptions{
		IncludeTimestamp: false,
	}
}

// SealOptions controls how a snapshot is sealed.
type SealOptions struct {
	// IncludeTimestamp includes the snapshot timestamp in the hash input.
	IncludeTimestamp bool
	// Clock overrides the current time source; defaults to time.Now.
	Clock func() time.Time
}

// Seal computes a deterministic SHA-256 hash over the snapshot's ports and
// returns a SealedSnapshot. The hash can later be verified with Verify.
func Seal(snap Snapshot, opts SealOptions) (SealedSnapshot, error) {
	if opts.Clock == nil {
		opts.Clock = time.Now
	}

	hashInput := struct {
		Ports     []Port    `json:"ports"`
		Timestamp time.Time `json:"timestamp,omitempty"`
	}{
		Ports: snap.Ports,
	}
	if opts.IncludeTimestamp {
		hashInput.Timestamp = snap.Timestamp
	}

	data, err := json.Marshal(hashInput)
	if err != nil {
		return SealedSnapshot{}, fmt.Errorf("seal: marshal: %w", err)
	}

	sum := sha256.Sum256(data)
	return SealedSnapshot{
		Snapshot: snap,
		Hash:     hex.EncodeToString(sum[:]),
		SealedAt: opts.Clock(),
	}, nil
}

// Verify recomputes the hash for the sealed snapshot and returns an error if
// it does not match the stored hash.
func Verify(sealed SealedSnapshot, opts SealOptions) error {
	resealed, err := Seal(sealed.Snapshot, opts)
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}
	if resealed.Hash != sealed.Hash {
		return errors.New("verify: hash mismatch — snapshot may have been tampered with")
	}
	return nil
}
