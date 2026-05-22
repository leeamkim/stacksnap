package diff

import (
	"fmt"
	"strings"

	"github.com/stacksnap/internal/snapshot"
)

// Result holds the comparison between two snapshots.
type Result struct {
	Added   []string
	Removed []string
	Changed []Change
	Equal   bool
}

// Change represents a tool or env var whose value changed between snapshots.
type Change struct {
	Name string
	From string
	To   string
}

// Compare computes the diff between a baseline and current snapshot.
func Compare(baseline, current *snapshot.Snapshot) *Result {
	r := &Result{}

	baseTools := toolMap(baseline.Tools)
	currTools := toolMap(current.Tools)

	for name, ver := range currTools {
		if bver, ok := baseTools[name]; !ok {
			r.Added = append(r.Added, fmt.Sprintf("%s@%s", name, ver))
		} else if bver != ver {
			r.Changed = append(r.Changed, Change{Name: name, From: bver, To: ver})
		}
	}

	for name, ver := range baseTools {
		if _, ok := currTools[name]; !ok {
			r.Removed = append(r.Removed, fmt.Sprintf("%s@%s", name, ver))
		}
	}

	r.Equal = len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0
	return r
}

// Summary returns a human-readable summary of the diff result.
func Summary(r *Result) string {
	if r.Equal {
		return "Snapshots are identical."
	}
	var sb strings.Builder
	for _, a := range r.Added {
		sb.WriteString(fmt.Sprintf("+ %s\n", a))
	}
	for _, rm := range r.Removed {
		sb.WriteString(fmt.Sprintf("- %s\n", rm))
	}
	for _, ch := range r.Changed {
		sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", ch.Name, ch.From, ch.To))
	}
	return sb.String()
}

func toolMap(tools []snapshot.Tool) map[string]string {
	m := make(map[string]string, len(tools))
	for _, t := range tools {
		m[t.Name] = t.Version
	}
	return m
}
