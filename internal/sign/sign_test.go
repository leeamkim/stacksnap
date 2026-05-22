package sign_test

import (
	"testing"

	"github.com/stacksnap/internal/sign"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID: "test-id-001",
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
		},
		Env: map[string]string{"GOPATH": "/home/user/go"},
		Annotations: map[string]string{"author": "alice"},
	}
}

func TestSign_AddsSignatureAnnotation(t *testing.T) {
	snap := makeSnap()
	if err := sign.Sign(snap, "supersecret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Annotations["stacksnap_sig"] == "" {
		t.Fatal("expected signature annotation to be set")
	}
}

func TestSign_NilSnapshotReturnsError(t *testing.T) {
	if err := sign.Sign(nil, "secret"); err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestSign_EmptySecretReturnsError(t *testing.T) {
	if err := sign.Sign(makeSnap(), ""); err == nil {
		t.Fatal("expected error for empty secret")
	}
}

func TestVerify_ValidSignature(t *testing.T) {
	snap := makeSnap()
	if err := sign.Sign(snap, "mysecret"); err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	if err := sign.Verify(snap, "mysecret"); err != nil {
		t.Fatalf("verify failed: %v", err)
	}
}

func TestVerify_WrongSecretFails(t *testing.T) {
	snap := makeSnap()
	if err := sign.Sign(snap, "correct"); err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	if err := sign.Verify(snap, "wrong"); err == nil {
		t.Fatal("expected verification to fail with wrong secret")
	}
}

func TestVerify_TamperedSnapshotFails(t *testing.T) {
	snap := makeSnap()
	if err := sign.Sign(snap, "secret"); err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	snap.Tools = append(snap.Tools, snapshot.Tool{Name: "node", Version: "20.0.0"})
	if err := sign.Verify(snap, "secret"); err == nil {
		t.Fatal("expected verification to fail after tampering")
	}
}

func TestVerify_MissingSignatureReturnsError(t *testing.T) {
	snap := makeSnap()
	if err := sign.Verify(snap, "secret"); err == nil {
		t.Fatal("expected error when no signature present")
	}
}

func TestVerify_NilSnapshotReturnsError(t *testing.T) {
	if err := sign.Verify(nil, "secret"); err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}
