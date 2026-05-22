// Package annotate provides functionality for adding and retrieving
// human-readable notes/annotations on snapshots.
package annotate

import (
	"errors"
	"strings"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

// Annotation holds a single note attached to a snapshot.
type Annotation struct {
	Author    string    `json:"author"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// Add appends an annotation to the snapshot's metadata map.
// Returns an error if note is empty or snap is nil.
func Add(snap *snapshot.Snapshot, author, note string) error {
	if snap == nil {
		return errors.New("annotate: snapshot must not be nil")
	}
	note = strings.TrimSpace(note)
	if note == "" {
		return errors.New("annotate: note must not be empty")
	}
	if snap.Meta == nil {
		snap.Meta = make(map[string]string)
	}
	author = strings.TrimSpace(author)
	if author == "" {
		author = "unknown"
	}
	snap.Meta["annotation.author"] = author
	snap.Meta["annotation.note"] = note
	snap.Meta["annotation.created_at"] = time.Now().UTC().Format(time.RFC3339)
	return nil
}

// Get retrieves the annotation stored in the snapshot's metadata.
// Returns nil if no annotation is present.
func Get(snap *snapshot.Snapshot) *Annotation {
	if snap == nil || snap.Meta == nil {
		return nil
	}
	note, ok := snap.Meta["annotation.note"]
	if !ok || note == "" {
		return nil
	}
	a := &Annotation{
		Author: snap.Meta["annotation.author"],
		Note:   note,
	}
	if ts, ok := snap.Meta["annotation.created_at"]; ok {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			a.CreatedAt = t
		}
	}
	return a
}

// Clear removes any annotation from the snapshot's metadata.
func Clear(snap *snapshot.Snapshot) error {
	if snap == nil {
		return errors.New("annotate: snapshot must not be nil")
	}
	if snap.Meta != nil {
		delete(snap.Meta, "annotation.author")
		delete(snap.Meta, "annotation.note")
		delete(snap.Meta, "annotation.created_at")
	}
	return nil
}
