package validate

import (
	"errors"
	"fmt"

	"github.com/stacksnap/internal/snapshot"
)

// Result holds the outcome of a validation check.
type Result struct {
	Tool    string
	Passed  bool
	Message string
}

// Report aggregates all validation results.
type Report struct {
	Results []Result
	Passed  int
	Failed  int
}

// Check validates that every tool in the snapshot has a non-empty version
// and that required tools are present.
func Check(snap snapshot.Snapshot, required []string) (Report, error) {
	if snap.Tools == nil {
		return Report{}, errors.New("snapshot contains no tools")
	}

	toolIndex := make(map[string]string, len(snap.Tools))
	for _, t := range snap.Tools {
		toolIndex[t.Name] = t.Version
	}

	var report Report

	// Check each detected tool has a version.
	for _, t := range snap.Tools {
		r := Result{Tool: t.Name}
		if t.Version == "" || t.Version == "unknown" {
			r.Passed = false
			r.Message = fmt.Sprintf("%s: version could not be determined", t.Name)
			report.Failed++
		} else {
			r.Passed = true
			r.Message = fmt.Sprintf("%s: %s", t.Name, t.Version)
			report.Passed++
		}
		report.Results = append(report.Results, r)
	}

	// Check required tools are present.
	for _, req := range required {
		if _, ok := toolIndex[req]; !ok {
			report.Results = append(report.Results, Result{
				Tool:    req,
				Passed:  false,
				Message: fmt.Sprintf("%s: required but not found", req),
			})
			report.Failed++
		}
	}

	return report, nil
}
