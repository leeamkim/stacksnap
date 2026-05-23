// Package checksum provides snapshot integrity verification via content hashing.
package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/stacksnap/internal/snapshot"
)

const AnnotationKey = "checksum.sha256"

// Compute returns a deterministic SHA-256 hex digest of the snapshot's
// tools and environment, excluding any existing checksum annotation.
func Compute(snap *snapshot.Snapshot) (string, error) {
	if snap == nil {
		return "", errors.New("checksum: snapshot is nil")
	}

	type stable struct {
		ID    string            `json:"id"`
		Tools []snapshot.Tool   `json:"tools"`
		Env   map[string]string `json:"env"`
	}

	// Sort tools by name for determinism.
	tools := make([]snapshot.Tool, len(snap.Tools))
	copy(tools, snap.Tools)
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})

	// Sort env keys for determinism.
	env := make(map[string]string, len(snap.Env))
	for k, v := range snap.Env {
		env[k] = v
	}

	s := stable{ID: snap.ID, Tools: tools, Env: env}
	data, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("checksum: marshal failed: %w", err)
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// Stamp computes the checksum and stores it in the snapshot's annotations.
func Stamp(snap *snapshot.Snapshot) error {
	digest, err := Compute(snap)
	if err != nil {
		return err
	}
	if snap.Annotations == nil {
		snap.Annotations = make(map[string]string)
	}
	snap.Annotations[AnnotationKey] = digest
	return nil
}

// Verify recomputes the checksum and compares it against the stored annotation.
// Returns an error if the annotation is missing or does not match.
func Verify(snap *snapshot.Snapshot) error {
	if snap == nil {
		return errors.New("checksum: snapshot is nil")
	}
	stored, ok := snap.Annotations[AnnotationKey]
	if !ok || stored == "" {
		return errors.New("checksum: no checksum annotation found")
	}
	expected, err := Compute(snap)
	if err != nil {
		return err
	}
	if stored != expected {
		return fmt.Errorf("checksum: mismatch — stored %s, computed %s", stored, expected)
	}
	return nil
}
