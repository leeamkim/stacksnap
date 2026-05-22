package diff_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stacksnap/internal/diff"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap(tools []snapshot.Tool) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		CapturedAt: time.Now(),
		Tools:      tools,
		Env:        map[string]string{},
	}
}

func TestCompare_Equal(t *testing.T) {
	tools := []snapshot.Tool{{Name: "go", Version: "1.22.0"}}
	r := diff.Compare(makeSnap(tools), makeSnap(tools))
	if !r.Equal {
		t.Error("expected snapshots to be equal")
	}
}

func TestCompare_Added(t *testing.T) {
	baseline := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.22.0"}})
	current := makeSnap([]snapshot.Tool{
		{Name: "go", Version: "1.22.0"},
		{Name: "node", Version: "20.0.0"},
	})
	r := diff.Compare(baseline, current)
	if len(r.Added) != 1 || !strings.Contains(r.Added[0], "node") {
		t.Errorf("expected node to be added, got %v", r.Added)
	}
}

func TestCompare_Removed(t *testing.T) {
	baseline := makeSnap([]snapshot.Tool{
		{Name: "go", Version: "1.22.0"},
		{Name: "node", Version: "20.0.0"},
	})
	current := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.22.0"}})
	r := diff.Compare(baseline, current)
	if len(r.Removed) != 1 || !strings.Contains(r.Removed[0], "node") {
		t.Errorf("expected node to be removed, got %v", r.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	baseline := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.21.0"}})
	current := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.22.0"}})
	r := diff.Compare(baseline, current)
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed tool, got %d", len(r.Changed))
	}
	if r.Changed[0].From != "1.21.0" || r.Changed[0].To != "1.22.0" {
		t.Errorf("unexpected change: %+v", r.Changed[0])
	}
}

func TestSummary_EqualMessage(t *testing.T) {
	r := &diff.Result{Equal: true}
	if !strings.Contains(diff.Summary(r), "identical") {
		t.Error("expected 'identical' in summary for equal result")
	}
}

func TestSummary_ContainsPrefixes(t *testing.T) {
	r := &diff.Result{
		Added:   []string{"node@20.0.0"},
		Removed: []string{"python@3.11"},
		Changed: []diff.Change{{Name: "go", From: "1.21", To: "1.22"}},
	}
	s := diff.Summary(r)
	if !strings.Contains(s, "+ node") {
		t.Error("expected '+' prefix for added")
	}
	if !strings.Contains(s, "- python") {
		t.Error("expected '-' prefix for removed")
	}
	if !strings.Contains(s, "~ go") {
		t.Error("expected '~' prefix for changed")
	}
}
