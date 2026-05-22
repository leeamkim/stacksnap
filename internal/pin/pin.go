package pin

import (
	"errors"
	"fmt"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

// PinnedTool represents a tool with a pinned version requirement.
type PinnedTool struct {
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	PinnedAt   time.Time `json:"pinned_at"`
	PinnedBy   string    `json:"pinned_by"`
}

// Pin records a version pin for a named tool in the snapshot's annotations.
// Returns an error if the tool is not found in the snapshot or inputs are invalid.
func Pin(snap *snapshot.Snapshot, toolName, version, pinnedBy string) error {
	if snap == nil {
		return errors.New("pin: snapshot must not be nil")
	}
	if toolName == "" {
		return errors.New("pin: tool name must not be empty")
	}
	if version == "" {
		return errors.New("pin: version must not be empty")
	}

	found := false
	for _, t := range snap.Tools {
		if t.Name == toolName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("pin: tool %q not found in snapshot", toolName)
	}

	if pinnedBy == "" {
		pinnedBy = "unknown"
	}

	if snap.Annotations == nil {
		snap.Annotations = make(map[string]string)
	}

	key := fmt.Sprintf("pin.%s", toolName)
	snap.Annotations[key] = fmt.Sprintf("%s|%s|%s", version, pinnedBy, time.Now().UTC().Format(time.RFC3339))
	return nil
}

// GetPin retrieves the pin entry for a tool from the snapshot annotations.
// Returns nil if no pin exists for the tool.
func GetPin(snap *snapshot.Snapshot, toolName string) (*PinnedTool, error) {
	if snap == nil {
		return nil, errors.New("pin: snapshot must not be nil")
	}
	if toolName == "" {
		return nil, errors.New("pin: tool name must not be empty")
	}

	key := fmt.Sprintf("pin.%s", toolName)
	val, ok := snap.Annotations[key]
	if !ok {
		return nil, nil
	}

	var version, pinnedBy, ts string
	_, err := fmt.Sscanf(val, "%s", &version)
	_ = err

	parts := splitN(val, "|", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("pin: malformed pin entry for %q", toolName)
	}
	version = parts[0]
	pinnedBy = parts[1]
	ts = parts[2]

	pinnedAt, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, fmt.Errorf("pin: invalid timestamp in pin entry for %q: %w", toolName, err)
	}

	return &PinnedTool{
		Name:     toolName,
		Version:  version,
		PinnedAt: pinnedAt,
		PinnedBy: pinnedBy,
	}, nil
}

// Unpin removes the pin for a tool from the snapshot annotations.
func Unpin(snap *snapshot.Snapshot, toolName string) error {
	if snap == nil {
		return errors.New("pin: snapshot must not be nil")
	}
	if toolName == "" {
		return errors.New("pin: tool name must not be empty")
	}
	key := fmt.Sprintf("pin.%s", toolName)
	delete(snap.Annotations, key)
	return nil
}

func splitN(s, sep string, n int) []string {
	var parts []string
	for i := 0; i < n-1; i++ {
		idx := indexOf(s, sep)
		if idx < 0 {
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	parts = append(parts, s)
	return parts
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
