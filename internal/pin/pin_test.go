package pin_test

import (
	"testing"

	"github.com/stacksnap/internal/pin"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID: "test-snap-001",
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
		},
		Annotations: make(map[string]string),
	}
}

func TestPin_AddsAnnotation(t *testing.T) {
	snap := makeSnap()
	err := pin.Pin(snap, "go", "1.22.0", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := snap.Annotations["pin.go"]; !ok {
		t.Error("expected annotation pin.go to be set")
	}
}

func TestPin_NilSnapshotReturnsError(t *testing.T) {
	err := pin.Pin(nil, "go", "1.22.0", "alice")
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestPin_UnknownToolReturnsError(t *testing.T) {
	snap := makeSnap()
	err := pin.Pin(snap, "rust", "1.76.0", "bob")
	if err == nil {
		t.Error("expected error for unknown tool")
	}
}

func TestPin_EmptyVersionReturnsError(t *testing.T) {
	snap := makeSnap()
	err := pin.Pin(snap, "go", "", "alice")
	if err == nil {
		t.Error("expected error for empty version")
	}
}

func TestGetPin_ReturnsPin(t *testing.T) {
	snap := makeSnap()
	if err := pin.Pin(snap, "node", "20.11.0", "carol"); err != nil {
		t.Fatalf("pin failed: %v", err)
	}
	p, err := pin.GetPin(snap, "node")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pin")
	}
	if p.Version != "20.11.0" {
		t.Errorf("expected version 20.11.0, got %s", p.Version)
	}
	if p.PinnedBy != "carol" {
		t.Errorf("expected pinnedBy carol, got %s", p.PinnedBy)
	}
}

func TestGetPin_MissingReturnsNil(t *testing.T) {
	snap := makeSnap()
	p, err := pin.GetPin(snap, "go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != nil {
		t.Error("expected nil for unpinned tool")
	}
}

func TestUnpin_RemovesAnnotation(t *testing.T) {
	snap := makeSnap()
	_ = pin.Pin(snap, "go", "1.22.0", "alice")
	if err := pin.Unpin(snap, "go"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := snap.Annotations["pin.go"]; ok {
		t.Error("expected pin.go annotation to be removed")
	}
}

func TestPin_EmptyPinnedByDefaultsToUnknown(t *testing.T) {
	snap := makeSnap()
	_ = pin.Pin(snap, "go", "1.22.0", "")
	p, err := pin.GetPin(snap, "go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.PinnedBy != "unknown" {
		t.Errorf("expected pinnedBy unknown, got %s", p.PinnedBy)
	}
}
