package filter

import (
	"strings"

	"github.com/stacksnap/internal/snapshot"
)

// Options controls which fields are included in the filtered snapshot.
type Options struct {
	Tools bool
	Env   bool
	OS    bool
	Tags  []string
}

// DefaultOptions returns an Options that includes everything.
func DefaultOptions() Options {
	return Options{
		Tools: true,
		Env:   true,
		OS:    true,
	}
}

// Apply returns a copy of snap with fields removed according to opts.
// If Tags is non-empty, only tools whose name contains at least one tag are kept.
func Apply(snap snapshot.Snapshot, opts Options) snapshot.Snapshot {
	out := snapshot.Snapshot{
		CapturedAt: snap.CapturedAt,
	}

	if opts.OS {
		out.OS = snap.OS
	}

	if opts.Env {
		out.Env = snap.Env
	}

	if opts.Tools {
		if len(opts.Tags) == 0 {
			out.Tools = snap.Tools
		} else {
			out.Tools = filterByTags(snap.Tools, opts.Tags)
		}
	}

	return out
}

// filterByTags returns only the tools whose name contains at least one of the
// provided tag strings (case-insensitive).
func filterByTags(tools []snapshot.Tool, tags []string) []snapshot.Tool {
	var result []snapshot.Tool
	for _, t := range tools {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(t.Name), strings.ToLower(tag)) {
				result = append(result, t)
				break
			}
		}
	}
	return result
}
