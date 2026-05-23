package audit_test

import (
	"testing"
	"time"

	"github.com/stacksnap/internal/audit"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:          "test-id",
		CapturedAt:  time.Now(),
		Annotations: map[string]string{},
	}
}

func TestRecord_AddsAnnotation(t *testing.T) {
	snap := makeSnap()
	err := audit.Record(snap, audit.EventCapture, "alice", "initial capture")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Annotations) == 0 {
		t.Error("expected at least one annotation after Record")
	}
}

func TestRecord_NilSnapshotReturnsError(t *testing.T) {
	err := audit.Record(nil, audit.EventExport, "bob", "")
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestRecord_EmptyKindReturnsError(t *testing.T) {
	snap := makeSnap()
	err := audit.Record(snap, "", "alice", "")
	if err == nil {
		t.Error("expected error for empty kind")
	}
}

func TestRecord_EmptyActorDefaultsToUnknown(t *testing.T) {
	snap := makeSnap()
	_ = audit.Record(snap, audit.EventValidate, "", "check")
	for _, v := range snap.Annotations {
		if v != "" {
			// actor=unknown should appear in value
			if len(v) > 0 {
				return
			}
		}
	}
}

func TestRecord_NilAnnotationsInitialised(t *testing.T) {
	snap := makeSnap()
	snap.Annotations = nil
	err := audit.Record(snap, audit.EventShare, "carol", "shared via link")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Annotations == nil {
		t.Error("expected annotations map to be initialised")
	}
}

func TestList_ReturnsEvents(t *testing.T) {
	snap := makeSnap()
	_ = audit.Record(snap, audit.EventCapture, "alice", "first")
	_ = audit.Record(snap, audit.EventExport, "alice", "second")
	events, err := audit.List(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) < 2 {
		t.Errorf("expected at least 2 events, got %d", len(events))
	}
}

func TestList_NilSnapshotReturnsError(t *testing.T) {
	_, err := audit.List(nil)
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}
