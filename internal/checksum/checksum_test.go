package checksum_test

import (
	"strings"
	"testing"

	"github.com/stacksnap/internal/checksum"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID: "test-id-001",
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
		},
		Env: map[string]string{
			"GOPATH": "/home/user/go",
			"NODE_ENV": "development",
		},
		Annotations: map[string]string{},
	}
}

func TestCompute_ReturnsHexString(t *testing.T) {
	snap := makeSnap()
	digest, err := checksum.Compute(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(digest) != 64 {
		t.Errorf("expected 64-char hex digest, got %d chars", len(digest))
	}
}

func TestCompute_NilSnapshotReturnsError(t *testing.T) {
	_, err := checksum.Compute(nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestCompute_IsDeterministic(t *testing.T) {
	snap := makeSnap()
	d1, _ := checksum.Compute(snap)
	d2, _ := checksum.Compute(snap)
	if d1 != d2 {
		t.Errorf("expected same digest on repeated calls, got %s vs %s", d1, d2)
	}
}

func TestCompute_ChangesWhenToolVersionChanges(t *testing.T) {
	s1 := makeSnap()
	s2 := makeSnap()
	s2.Tools[0].Version = "1.21.0"

	d1, _ := checksum.Compute(s1)
	d2, _ := checksum.Compute(s2)
	if d1 == d2 {
		t.Error("expected different digests for different tool versions")
	}
}

func TestStamp_AddsAnnotation(t *testing.T) {
	snap := makeSnap()
	if err := checksum.Stamp(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := snap.Annotations[checksum.AnnotationKey]
	if !ok {
		t.Fatal("expected annotation to be set")
	}
	if !strings.HasPrefix(val, "") || len(val) != 64 {
		t.Errorf("unexpected annotation value: %s", val)
	}
}

func TestVerify_ValidSnapshot(t *testing.T) {
	snap := makeSnap()
	_ = checksum.Stamp(snap)
	if err := checksum.Verify(snap); err != nil {
		t.Errorf("expected valid snapshot to pass verification: %v", err)
	}
}

func TestVerify_TamperedSnapshot(t *testing.T) {
	snap := makeSnap()
	_ = checksum.Stamp(snap)
	snap.Tools[0].Version = "9.9.9" // tamper
	if err := checksum.Verify(snap); err == nil {
		t.Error("expected verification to fail after tampering")
	}
}

func TestVerify_MissingAnnotationReturnsError(t *testing.T) {
	snap := makeSnap()
	if err := checksum.Verify(snap); err == nil {
		t.Error("expected error when no checksum annotation present")
	}
}
