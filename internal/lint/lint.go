package lint

import (
	"fmt"
	"strings"

	"github.com/stacksnap/internal/snapshot"
)

// Issue represents a single lint warning or error found in a snapshot.
type Issue struct {
	Level   string // "warn" or "error"
	Field   string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(i.Level), i.Field, i.Message)
}

// Result holds all issues found during a lint run.
type Result struct {
	Issues []Issue
}

// OK returns true if there are no error-level issues.
func (r *Result) OK() bool {
	for _, iss := range r.Issues {
		if iss.Level == "error" {
			return false
		}
	}
	return true
}

// Check runs all lint rules against the given snapshot and returns a Result.
func Check(snap *snapshot.Snapshot) (*Result, error) {
	if snap == nil {
		return nil, fmt.Errorf("lint: snapshot must not be nil")
	}

	var issues []Issue

	// Rule: snapshot must have an ID
	if strings.TrimSpace(snap.ID) == "" {
		issues = append(issues, Issue{Level: "error", Field: "id", Message: "snapshot ID is empty"})
	}

	// Rule: at least one tool should be present
	if len(snap.Tools) == 0 {
		issues = append(issues, Issue{Level: "warn", Field: "tools", Message: "no tools detected in snapshot"})
	}

	// Rule: tool versions should not be "unknown"
	for _, t := range snap.Tools {
		if strings.EqualFold(t.Version, "unknown") || strings.TrimSpace(t.Version) == "" {
			issues = append(issues, Issue{
				Level:   "warn",
				Field:   fmt.Sprintf("tools[%s].version", t.Name),
				Message: "version is unknown or empty",
			})
		}
	}

	// Rule: env keys should not be blank
	for k := range snap.Env {
		if strings.TrimSpace(k) == "" {
			issues = append(issues, Issue{Level: "error", Field: "env", Message: "env contains a blank key"})
			break
		}
	}

	return &Result{Issues: issues}, nil
}
