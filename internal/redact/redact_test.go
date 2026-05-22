package redact_test

import (
	"testing"

	"github.com/stacksnap/internal/redact"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap(env map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID: "test-id",
		Env: env,
	}
}

func TestApply_NilSnapshotReturnsError(t *testing.T) {
	_, err := redact.Apply(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	snap := makeSnap(map[string]string{"API_KEY": "secret123", "HOME": "/home/user"})
	_, err := redact.Apply(snap, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Env["API_KEY"] != "secret123" {
		t.Error("original snapshot was mutated")
	}
}

func TestApply_RedactsSensitiveKeys(t *testing.T) {
	snap := makeSnap(map[string]string{
		"API_KEY":  "abc",
		"PASSWORD": "hunter2",
		"HOME":     "/home/user",
		"TOKEN":    "tok123",
	})
	result, err := redact.Apply(snap, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"API_KEY", "PASSWORD", "TOKEN"} {
		if result.Env[key] != "[REDACTED]" {
			t.Errorf("expected %s to be redacted, got %q", key, result.Env[key])
		}
	}
	if result.Env["HOME"] != "/home/user" {
		t.Errorf("expected HOME to be preserved, got %q", result.Env["HOME"])
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	snap := makeSnap(map[string]string{"SECRET_KEY": "s3cr3t"})
	opts := &redact.Options{
		Patterns:    redact.DefaultPatterns,
		Placeholder: "***",
	}
	result, err := redact.Apply(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["SECRET_KEY"] != "***" {
		t.Errorf("expected custom placeholder, got %q", result.Env["SECRET_KEY"])
	}
}

func TestApply_InvalidPatternReturnsError(t *testing.T) {
	snap := makeSnap(map[string]string{"FOO": "bar"})
	opts := &redact.Options{Patterns: []string{"[invalid"}}
	_, err := redact.Apply(snap, opts)
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestSensitiveKeys_ReturnsSortedList(t *testing.T) {
	snap := makeSnap(map[string]string{
		"TOKEN":    "t",
		"API_KEY":  "k",
		"HOSTNAME": "h",
	})
	keys := redact.SensitiveKeys(snap, nil)
	if len(keys) != 2 {
		t.Fatalf("expected 2 sensitive keys, got %d: %v", len(keys), keys)
	}
	if keys[0] > keys[1] {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}

func TestSensitiveKeys_NilSnapshotReturnsNil(t *testing.T) {
	keys := redact.SensitiveKeys(nil, nil)
	if keys != nil {
		t.Errorf("expected nil for nil snapshot, got %v", keys)
	}
}
