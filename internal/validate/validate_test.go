package validate_test

import (
	"testing"

	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/validate"
)

func makeSnap(tools []snapshot.Tool) snapshot.Snapshot {
	return snapshot.Snapshot{
		Tools: tools,
		Env:   map[string]string{},
	}
}

func TestCheck_AllToolsValid(t *testing.T) {
	snap := makeSnap([]snapshot.Tool{
		{Name: "go", Version: "1.22.0"},
		{Name: "node", Version: "20.11.0"},
	})
	report, err := validate.Check(snap, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Failed != 0 {
		t.Errorf("expected 0 failures, got %d", report.Failed)
	}
	if report.Passed != 2 {
		t.Errorf("expected 2 passed, got %d", report.Passed)
	}
}

func TestCheck_UnknownVersionFails(t *testing.T) {
	snap := makeSnap([]snapshot.Tool{
		{Name: "go", Version: "1.22.0"},
		{Name: "rust", Version: "unknown"},
	})
	report, err := validate.Check(snap, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Failed != 1 {
		t.Errorf("expected 1 failure, got %d", report.Failed)
	}
}

func TestCheck_RequiredToolMissing(t *testing.T) {
	snap := makeSnap([]snapshot.Tool{
		{Name: "go", Version: "1.22.0"},
	})
	report, err := validate.Check(snap, []string{"docker"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Failed != 1 {
		t.Errorf("expected 1 failure for missing required tool, got %d", report.Failed)
	}
}

func TestCheck_NilToolsReturnsError(t *testing.T) {
	snap := snapshot.Snapshot{}
	_, err := validate.Check(snap, nil)
	if err == nil {
		t.Error("expected error for nil tools, got nil")
	}
}

func TestCheck_ResultMessagesPopulated(t *testing.T) {
	snap := makeSnap([]snapshot.Tool{
		{Name: "go", Version: "1.22.0"},
	})
	report, _ := validate.Check(snap, nil)
	for _, r := range report.Results {
		if r.Message == "" {
			t.Errorf("expected non-empty message for tool %s", r.Tool)
		}
	}
}
