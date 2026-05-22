// Package share provides functionality to encode and decode snapshots
// as shareable URL-safe strings (base64-encoded JSON).
package share

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/user/stacksnap/internal/snapshot"
)

const SchemePrefix = "stacksnap://"

// Encode serializes a snapshot to a URL-safe base64 string prefixed
// with the stacksnap:// scheme.
func Encode(snap snapshot.Snapshot) (string, error) {
	data, err := json.Marshal(snap)
	if err != nil {
		return "", fmt.Errorf("share: marshal snapshot: %w", err)
	}
	encoded := base64.URLEncoding.EncodeToString(data)
	return SchemePrefix + encoded, nil
}

// Decode parses a share string (with or without the scheme prefix)
// and returns the embedded snapshot.
func Decode(shareStr string) (snapshot.Snapshot, error) {
	var snap snapshot.Snapshot

	payload := shareStr
	if len(shareStr) > len(SchemePrefix) && shareStr[:len(SchemePrefix)] == SchemePrefix {
		payload = shareStr[len(SchemePrefix):]
	}

	data, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return snap, fmt.Errorf("share: base64 decode: %w", err)
	}

	if err := json.Unmarshal(data, &snap); err != nil {
		return snap, fmt.Errorf("share: unmarshal snapshot: %w", err)
	}

	return snap, nil
}
