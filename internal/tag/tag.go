package tag

import (
	"errors"
	"regexp"
	"strings"

	"github.com/stacksnap/internal/snapshot"
)

var validTagRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// ErrInvalidTag is returned when a tag name is malformed.
var ErrInvalidTag = errors.New("tag: invalid tag name (alphanumeric, dash, underscore only)")

// Add appends one or more tags to the snapshot, deduplicating as it goes.
// Returns ErrInvalidTag if any tag name is malformed.
func Add(snap *snapshot.Snapshot, tags ...string) error {
	for _, t := range tags {
		if !validTagRe.MatchString(t) {
			return ErrInvalidTag
		}
	}
	existing := make(map[string]struct{}, len(snap.Tags))
	for _, t := range snap.Tags {
		existing[t] = struct{}{}
	}
	for _, t := range tags {
		if _, ok := existing[t]; !ok {
			snap.Tags = append(snap.Tags, t)
			existing[t] = struct{}{}
		}
	}
	return nil
}

// Remove deletes the given tags from the snapshot. Unknown tags are silently
// ignored.
func Remove(snap *snapshot.Snapshot, tags ...string) {
	remove := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		remove[t] = struct{}{}
	}
	kept := snap.Tags[:0]
	for _, t := range snap.Tags {
		if _, skip := remove[t]; !skip {
			kept = append(kept, t)
		}
	}
	snap.Tags = kept
}

// List returns a sorted, deduplicated copy of the snapshot's tags.
func List(snap *snapshot.Snapshot) []string {
	seen := make(map[string]struct{}, len(snap.Tags))
	out := make([]string, 0, len(snap.Tags))
	for _, t := range snap.Tags {
		if _, ok := seen[t]; !ok {
			out = append(out, t)
			seen[t] = struct{}{}
		}
	}
	// simple insertion sort – tag lists are tiny
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && strings.ToLower(out[j]) < strings.ToLower(out[j-1]); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
