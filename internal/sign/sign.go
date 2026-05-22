package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stacksnap/internal/snapshot"
)

const signatureKey = "stacksnap_sig"

// Sign computes an HMAC-SHA256 signature over the snapshot's canonical JSON
// representation and stores it in the snapshot's Annotations map under
// "stacksnap_sig". The secret must be non-empty.
func Sign(snap *snapshot.Snapshot, secret string) error {
	if snap == nil {
		return errors.New("sign: snapshot is nil")
	}
	if secret == "" {
		return errors.New("sign: secret must not be empty")
	}

	sig, err := computeSignature(snap, secret)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	if snap.Annotations == nil {
		snap.Annotations = make(map[string]string)
	}
	snap.Annotations[signatureKey] = sig
	return nil
}

// Verify checks that the signature stored in the snapshot's Annotations
// matches a freshly computed HMAC-SHA256 over the snapshot (excluding the
// stored signature itself). Returns nil on success.
func Verify(snap *snapshot.Snapshot, secret string) error {
	if snap == nil {
		return errors.New("verify: snapshot is nil")
	}
	if secret == "" {
		return errors.New("verify: secret must not be empty")
	}

	stored, ok := snap.Annotations[signatureKey]
	if !ok || stored == "" {
		return errors.New("verify: no signature found in snapshot")
	}

	// Temporarily remove signature before computing digest.
	delete(snap.Annotations, signatureKey)
	expected, err := computeSignature(snap, secret)
	snap.Annotations[signatureKey] = stored // restore
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}

	if !hmac.Equal([]byte(stored), []byte(expected)) {
		return errors.New("verify: signature mismatch")
	}
	return nil
}

func computeSignature(snap *snapshot.Snapshot, secret string) (string, error) {
	data, err := json.Marshal(snap)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil)), nil
}
