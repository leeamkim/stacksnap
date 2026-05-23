package score_test

import (
	"testing"

	"github.com/stacksnap/internal/score"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID: "snap-001",
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
		},
		Env:  map[string]string{"GOPATH": "/home/user/go"},
		Tags: []string{"backend"},
		Annotations: map[string]string{"description": "my dev stack"},
	}
}

func TestEvaluate_PerfectSnapshot(t *testing.T) {
	res, err := score.Evaluate(makeSnap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 100 {
		t.Errorf("expected score 100, got %d; reasons: %v", res.Score, res.Reasons)
	}
	if res.Grade != score.GradeA {
		t.Errorf("expected grade A, got %s", res.Grade)
	}
}

func TestEvaluate_NilSnapshotReturnsError(t *testing.T) {
	_, err := score.Evaluate(nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestEvaluate_MissingIDDeducts(t *testing.T) {
	snap := makeSnap()
	snap.ID = ""
	res, err := score.Evaluate(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 85 {
		t.Errorf("expected score 85, got %d", res.Score)
	}
}

func TestEvaluate_UnknownToolVersionDeducts(t *testing.T) {
	snap := makeSnap()
	snap.Tools[0].Version = "unknown"
	res, err := score.Evaluate(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 95 {
		t.Errorf("expected score 95, got %d", res.Score)
	}
}

func TestEvaluate_GradeF_WhenScoreLow(t *testing.T) {
	snap := &snapshot.Snapshot{} // missing everything
	res, err := score.Evaluate(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Grade != score.GradeF {
		t.Errorf("expected grade F for empty snapshot, got %s (score=%d)", res.Grade, res.Score)
	}
	if len(res.Reasons) == 0 {
		t.Error("expected at least one reason for low score")
	}
}

func TestEvaluate_ScoreNeverBelowZero(t *testing.T) {
	snap := &snapshot.Snapshot{}
	res, _ := score.Evaluate(snap)
	if res.Score < 0 {
		t.Errorf("score must not be negative, got %d", res.Score)
	}
}
