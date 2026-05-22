package lint_test

import (
	"testing"

	"github.com/stacksnap/internal/lint"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID: "snap-abc123",
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
		},
		Env: map[string]string{
			"GOPATH": "/home/user/go",
		},
	}
}

func TestCheck_CleanSnapshotHasNoIssues(t *testing.T) {
	result, err := lint.Check(makeSnap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
	if !result.OK() {
		t.Error("expected OK() to be true")
	}
}

func TestCheck_MissingIDIsError(t *testing.T) {
	snap := makeSnap()
	snap.ID = ""
	result, err := lint.Check(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OK() {
		t.Error("expected OK() to be false when ID is missing")
	}
	if len(result.Issues) == 0 {
		t.Error("expected at least one issue")
	}
}

func TestCheck_NoToolsIsWarn(t *testing.T) {
	snap := makeSnap()
	snap.Tools = nil
	result, err := lint.Check(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK() {
		t.Error("expected OK() true (warn only)")
	}
	if len(result.Issues) == 0 {
		t.Error("expected a warn issue for missing tools")
	}
}

func TestCheck_UnknownVersionIsWarn(t *testing.T) {
	snap := makeSnap()
	snap.Tools = []snapshot.Tool{{Name: "rust", Version: "unknown"}}
	result, _ := lint.Check(snap)
	if len(result.Issues) == 0 {
		t.Error("expected warn for unknown version")
	}
	if result.Issues[0].Level != "warn" {
		t.Errorf("expected warn level, got %s", result.Issues[0].Level)
	}
}

func TestCheck_NilSnapshotReturnsError(t *testing.T) {
	_, err := lint.Check(nil)
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}
