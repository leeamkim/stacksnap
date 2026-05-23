package audit_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stacksnap/internal/audit"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:          "test-snap-001",
		CapturedAt:  time.Now(),
		Tools:       []snapshot.Tool{{Name: "go", Version: "1.21.0"}},
		Env:         map[string]string{"GOPATH": "/home/user/go"},
		Annotations: map[string]string{},
	}
}

func TestRecord_AddsAnnotation(t *testing.T) {
	snap := makeSnap()
	err := audit.Record(snap, audit.Entry{
		Kind:  "export",
		Actor: "alice",
		Note:  "exported to JSON",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := audit.List(snap)
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Kind != "export" {
		t.Errorf("expected kind 'export', got %q", entries[0].Kind)
	}
	if entries[0].Actor != "alice" {
		t.Errorf("expected actor 'alice', got %q", entries[0].Actor)
	}
}

func TestRecord_NilSnapshotReturnsError(t *testing.T) {
	err := audit.Record(nil, audit.Entry{Kind: "export", Actor: "alice"})
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestRecord_EmptyKindReturnsError(t *testing.T) {
	snap := makeSnap()
	err := audit.Record(snap, audit.Entry{Kind: "", Actor: "alice"})
	if err == nil {
		t.Fatal("expected error for empty kind")
	}
}

func TestRecord_EmptyActorDefaultsToUnknown(t *testing.T) {
	snap := makeSnap()
	err := audit.Record(snap, audit.Entry{Kind: "restore", Actor: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := audit.List(snap)
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if entries[0].Actor != "unknown" {
		t.Errorf("expected actor 'unknown', got %q", entries[0].Actor)
	}
}

func TestRecord_TimestampIsSet(t *testing.T) {
	snap := makeSnap()
	before := time.Now()
	_ = audit.Record(snap, audit.Entry{Kind: "validate", Actor: "ci"})
	after := time.Now()

	entries, _ := audit.List(snap)
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	ts := entries[0].At
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestList_MultipleEntriesOrdered(t *testing.T) {
	snap := makeSnap()
	_ = audit.Record(snap, audit.Entry{Kind: "export", Actor: "alice"})
	time.Sleep(2 * time.Millisecond)
	_ = audit.Record(snap, audit.Entry{Kind: "restore", Actor: "bob"})

	entries, err := audit.List(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// Most recent last
	if !entries[0].At.Before(entries[1].At) && entries[0].At != entries[1].At {
		t.Errorf("expected entries ordered oldest-first")
	}
}

func TestList_NilSnapshotReturnsError(t *testing.T) {
	_, err := audit.List(nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestList_EmptyAnnotationsReturnsEmpty(t *testing.T) {
	snap := makeSnap()
	entries, err := audit.List(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_NoteIsPreserved(t *testing.T) {
	snap := makeSnap()
	_ = audit.Record(snap, audit.Entry{Kind: "sign", Actor: "admin", Note: "signed with key abc"})

	entries, _ := audit.List(snap)
	if !strings.Contains(entries[0].Note, "signed with key abc") {
		t.Errorf("expected note to be preserved, got %q", entries[0].Note)
	}
}
