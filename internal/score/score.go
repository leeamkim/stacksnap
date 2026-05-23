package score

import (
	"fmt"
	"strings"

	"github.com/stacksnap/internal/snapshot"
)

// Grade represents a letter grade for a snapshot.
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeF Grade = "F"
)

// Result holds the numeric score and derived grade for a snapshot.
type Result struct {
	Score   int    `json:"score"`
	Grade   Grade  `json:"grade"`
	Reasons []string `json:"reasons"`
}

// Evaluate scores a snapshot on a 0–100 scale based on completeness and
// quality heuristics. Each deduction is recorded in Reasons.
func Evaluate(snap *snapshot.Snapshot) (*Result, error) {
	if snap == nil {
		return nil, fmt.Errorf("score: snapshot must not be nil")
	}

	score := 100
	var reasons []string

	deduct := func(pts int, reason string) {
		score -= pts
		reasons = append(reasons, reason)
	}

	if snap.ID == "" {
		deduct(15, "missing snapshot ID")
	}

	if len(snap.Tools) == 0 {
		deduct(20, "no tools detected")
	} else {
		for _, t := range snap.Tools {
			if t.Version == "" || strings.EqualFold(t.Version, "unknown") {
				deduct(5, fmt.Sprintf("tool %q has unknown version", t.Name))
			}
		}
	}

	if len(snap.Env) == 0 {
		deduct(10, "no environment variables captured")
	}

	if len(snap.Tags) == 0 {
		deduct(5, "no tags applied")
	}

	if _, ok := snap.Annotations["description"]; !ok {
		deduct(5, "missing description annotation")
	}

	if score < 0 {
		score = 0
	}

	return &Result{
		Score:   score,
		Grade:   gradeFor(score),
		Reasons: reasons,
	}, nil
}

func gradeFor(score int) Grade {
	switch {
	case score >= 90:
		return GradeA
	case score >= 75:
		return GradeB
	case score >= 55:
		return GradeC
	default:
		return GradeF
	}
}
