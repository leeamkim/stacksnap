package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

// Entry represents a single historical snapshot record.
type Entry struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label,omitempty"`
	Snapshot  snapshot.Snapshot `json:"snapshot"`
}

// DefaultDir is the default directory used to store history entries.
const DefaultDir = ".stacksnap/history"

// Save persists a snapshot as a history entry under dir.
// If label is empty, the entry ID is used as the label.
func Save(snap snapshot.Snapshot, label, dir string) (Entry, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Entry{}, fmt.Errorf("history: create dir: %w", err)
	}

	now := time.Now().UTC()
	id := now.Format("20060102T150405Z")
	if label == "" {
		label = id
	}

	entry := Entry{
		ID:        id,
		Timestamp: now,
		Label:     label,
		Snapshot:  snap,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return Entry{}, fmt.Errorf("history: marshal: %w", err)
	}

	path := filepath.Join(dir, id+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return Entry{}, fmt.Errorf("history: write file: %w", err)
	}

	return entry, nil
}

// List returns all history entries in dir sorted by timestamp ascending.
func List(dir string) ([]Entry, error) {
	glob := filepath.Join(dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}

	var entries []Entry
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("history: read %s: %w", path, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: unmarshal %s: %w", path, err)
		}
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	return entries, nil
}
