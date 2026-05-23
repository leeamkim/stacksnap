package compare_test

import (
	"testing"

	"github.com/stacksnap/internal/compare"
	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/template"
)

func makeSnap(tools map[string]string) *snapshot.Snapshot {
	s := &snapshot.Snapshot{ID: "test-snap"}
	for name, ver := range tools {
		s.Tools = append(s.Tools, snapshot.Tool{Name: name, Version: ver})
	}
	return s
}

func makeTmpl(tools map[string]string) *template.Template {
	t := &template.Template{Name: "base"}
	for name, ver := range tools {
		t.Tools = append(t.Tools, snapshot.Tool{Name: name, Version: ver})
	}
	return t
}

func TestAgainstTemplate_NilSnapshotReturnsError(t *testing.T) {
	_, err := compare.AgainstTemplate(nil, makeTmpl(nil), nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestAgainstTemplate_NilTemplateReturnsError(t *testing.T) {
	_, err := compare.AgainstTemplate(makeSnap(nil), nil, nil)
	if err == nil {
		t.Fatal("expected error for nil template")
	}
}

func TestAgainstTemplate_FullMatch(t *testing.T) {
	snap := makeSnap(map[string]string{"go": "1.22.0", "node": "20.0.0"})
	tmpl := makeTmpl(map[string]string{"go": "1.22.0", "node": "20.0.0"})
	res, err := compare.AgainstTemplate(snap, tmpl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK() {
		t.Errorf("expected OK, got %s", res.Summary())
	}
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched))
	}
}

func TestAgainstTemplate_MissingTool(t *testing.T) {
	snap := makeSnap(map[string]string{"go": "1.22.0"})
	tmpl := makeTmpl(map[string]string{"go": "1.22.0", "node": "20.0.0"})
	res, err := compare.AgainstTemplate(snap, tmpl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OK() {
		t.Error("expected not OK due to missing tool")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "node" {
		t.Errorf("expected missing=[node], got %v", res.Missing)
	}
}

func TestAgainstTemplate_VersionMismatch(t *testing.T) {
	snap := makeSnap(map[string]string{"go": "1.21.0"})
	tmpl := makeTmpl(map[string]string{"go": "1.22.0"})
	res, err := compare.AgainstTemplate(snap, tmpl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OK() {
		t.Error("expected not OK due to version mismatch")
	}
	if len(res.Mismatch) != 1 {
		t.Errorf("expected 1 mismatch, got %d", len(res.Mismatch))
	}
}

func TestAgainstTemplate_ExtraToolsAreRecorded(t *testing.T) {
	snap := makeSnap(map[string]string{"go": "1.22.0", "rust": "1.78.0"})
	tmpl := makeTmpl(map[string]string{"go": "1.22.0"})
	res, err := compare.AgainstTemplate(snap, tmpl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK() {
		t.Errorf("extra tools should not fail OK: %s", res.Summary())
	}
	if len(res.Extra) != 1 || res.Extra[0] != "rust" {
		t.Errorf("expected extra=[rust], got %v", res.Extra)
	}
}
