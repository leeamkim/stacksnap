// Package compare provides snapshot comparison against a baseline template.
package compare

import (
	"errors"
	"fmt"
	"strings"

	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/template"
)

// Result holds the outcome of comparing a snapshot against a template baseline.
type Result struct {
	Matched  []string
	Missing  []string
	Extra    []string
	Mismatch []string
}

// OK returns true when the snapshot fully satisfies the template baseline.
func (r *Result) OK() bool {
	return len(r.Missing) == 0 && len(r.Mismatch) == 0
}

// Summary returns a human-readable summary string.
func (r *Result) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "matched=%d missing=%d extra=%d mismatch=%d",
		len(r.Matched), len(r.Missing), len(r.Extra), len(r.Mismatch))
	return sb.String()
}

// AgainstTemplate compares snap against the tools declared in tmpl.
// Variables in the template are expanded using vars before comparison.
func AgainstTemplate(snap *snapshot.Snapshot, tmpl *template.Template, vars map[string]string) (*Result, error) {
	if snap == nil {
		return nil, errors.New("compare: snapshot must not be nil")
	}
	if tmpl == nil {
		return nil, errors.New("compare: template must not be nil")
	}

	instantiated, err := template.Instantiate(tmpl, vars)
	if err != nil {
		return nil, fmt.Errorf("compare: instantiate template: %w", err)
	}

	snapIndex := make(map[string]string, len(snap.Tools))
	for _, t := range snap.Tools {
		snapIndex[t.Name] = t.Version
	}

	tmplIndex := make(map[string]string, len(instantiated.Tools))
	for _, t := range instantiated.Tools {
		tmplIndex[t.Name] = t.Version
	}

	res := &Result{}

	for name, wantVer := range tmplIndex {
		gotVer, found := snapIndex[name]
		if !found {
			res.Missing = append(res.Missing, name)
			continue
		}
		if wantVer != "" && gotVer != wantVer {
			res.Mismatch = append(res.Mismatch, fmt.Sprintf("%s: want %s got %s", name, wantVer, gotVer))
		} else {
			res.Matched = append(res.Matched, name)
		}
	}

	for name := range snapIndex {
		if _, inTmpl := tmplIndex[name]; !inTmpl {
			res.Extra = append(res.Extra, name)
		}
	}

	return res, nil
}
