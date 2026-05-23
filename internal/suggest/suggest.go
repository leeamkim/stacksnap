// Package suggest analyzes a snapshot and recommends improvements or missing
// tools based on what has already been detected in the dev stack.
package suggest

import (
	"fmt"
	"strings"

	"github.com/yourusername/stacksnap/internal/snapshot"
)

// Severity indicates how strongly a suggestion should be acted upon.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Suggestion represents a single actionable recommendation.
type Suggestion struct {
	Tool     string   `json:"tool"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
}

// rule is an internal check applied to a snapshot.
type rule struct {
	name  string
	check func(snap *snapshot.Snapshot) []Suggestion
}

// builtinRules contains the default set of suggestion rules.
var builtinRules = []rule{
	{
		name: "missing-formatter",
		check: func(snap *snapshot.Snapshot) []Suggestion {
			formatters := []string{"prettier", "gofmt", "black", "rustfmt"}
			for _, t := range snap.Tools {
				for _, f := range formatters {
					if strings.EqualFold(t.Name, f) {
						return nil
					}
				}
			}
			return []Suggestion{{
				Tool:     "formatter",
				Message:  "No code formatter detected (e.g. prettier, gofmt, black). Consider adding one for consistent style.",
				Severity: SeverityWarning,
			}}
		},
	},
	{
		name: "missing-linter",
		check: func(snap *snapshot.Snapshot) []Suggestion {
			linters := []string{"eslint", "golangci-lint", "flake8", "clippy", "rubocop"}
			for _, t := range snap.Tools {
				for _, l := range linters {
					if strings.EqualFold(t.Name, l) {
						return nil
					}
				}
			}
			return []Suggestion{{
				Tool:     "linter",
				Message:  "No linter detected (e.g. eslint, golangci-lint, flake8). Adding a linter improves code quality.",
				Severity: SeverityWarning,
			}}
		},
	},
	{
		name: "missing-container-runtime",
		check: func(snap *snapshot.Snapshot) []Suggestion {
			runtimes := []string{"docker", "podman", "nerdctl"}
			for _, t := range snap.Tools {
				for _, r := range runtimes {
					if strings.EqualFold(t.Name, r) {
						return nil
					}
				}
			}
			return []Suggestion{{
				Tool:     "container-runtime",
				Message:  "No container runtime detected (e.g. docker, podman). Containerisation aids reproducibility.",
				Severity: SeverityInfo,
			}}
		},
	},
	{
		name: "unknown-tool-versions",
		check: func(snap *snapshot.Snapshot) []Suggestion {
			var suggestions []Suggestion
			for _, t := range snap.Tools {
				if t.Version == "unknown" || t.Version == "" {
					suggestions = append(suggestions, Suggestion{
						Tool:    t.Name,
						Message: fmt.Sprintf("Version of %q could not be determined. Pin an explicit version for reproducibility.", t.Name),
						Severity: SeverityCritical,
					})
				}
			}
			return suggestions
		},
	},
}

// Analyze runs all built-in suggestion rules against the provided snapshot and
// returns a slice of actionable Suggestions. An error is returned only if the
// snapshot itself is nil.
func Analyze(snap *snapshot.Snapshot) ([]Suggestion, error) {
	if snap == nil {
		return nil, fmt.Errorf("suggest: snapshot must not be nil")
	}

	var results []Suggestion
	for _, r := range builtinRules {
		results = append(results, r.check(snap)...)
	}
	return results, nil
}
