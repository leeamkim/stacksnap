package annotate_test

import (
	"testing"
	"time"

	"github.com/stacksnap/internal/annotate"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        "test-id",
		CreatedAt: time.Now(),
		Meta:      make(map[string]string),
	}
}

func TestAdd_StoresAnnotation(t *testing.T) {
	snap := makeSnap()
	if err := annotate.Add(snap, "alice", "initial setup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := annotate.Get(snap)
	if a == nil {
		t.Fatal("expected annotation, got nil")
	}
	if a.Author != "alice" {
		t.Errorf("expected author 'alice', got %q", a.Author)
	}
	if a.Note != "initial setup" {
		t.Errorf("expected note 'initial setup', got %q", a.Note)
	}
}

func TestAdd_EmptyNoteReturnsError(t *testing.T) {
	snap := makeSnap()
	if err := annotate.Add(snap, "alice", "   "); err == nil {
		t.Error("expected error for empty note, got nil")
	}
}

func TestAdd_NilSnapshotReturnsError(t *testing.T) {
	if err := annotate.Add(nil, "alice", "note"); err == nil {
		t.Error("expected error for nil snapshot, got nil")
	}
}

func TestAdd_EmptyAuthorDefaultsToUnknown(t *testing.T) {
	snap := makeSnap()
	if err := annotate.Add(snap, "", "some note"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := annotate.Get(snap)
	if a.Author != "unknown" {
		t.Errorf("expected 'unknown', got %q", a.Author)
	}
}

func TestGet_NilSnapshotReturnsNil(t *testing.T) {
	if annotate.Get(nil) != nil {
		t.Error("expected nil for nil snapshot")
	}
}

func TestGet_NoAnnotationReturnsNil(t *testing.T) {
	snap := makeSnap()
	if annotate.Get(snap) != nil {
		t.Error("expected nil when no annotation set")
	}
}

func TestClear_RemovesAnnotation(t *testing.T) {
	snap := makeSnap()
	_ = annotate.Add(snap, "bob", "clean up later")
	if err := annotate.Clear(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if annotate.Get(snap) != nil {
		t.Error("expected nil after clear")
	}
}

func TestClear_NilSnapshotReturnsError(t *testing.T) {
	if err := annotate.Clear(nil); err == nil {
		t.Error("expected error for nil snapshot")
	}
}
