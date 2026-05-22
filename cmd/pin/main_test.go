package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "pin")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binPath
}

func writeJSON(t *testing.T, snap *snapshot.Snapshot) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("failed to marshal snapshot: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write snapshot file: %v", err)
	}
	return path
}

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID: "test-snap-001",
		Timestamp: time.Now(),
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.21.0"},
			{Name: "node", Version: "20.0.0"},
		},
		Env: map[string]string{"GOPATH": "/home/user/go"},
	}
}

func TestMain_BuildsBinary(t *testing.T) {
	buildBinary(t)
}

func TestMain_MissingSnapshotFlagExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--tool", "go", "--version", "1.22.0")
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit when --snapshot flag is missing")
	}
}

func TestMain_MissingToolFlagExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	snap := makeSnap()
	path := writeJSON(t, snap)

	cmd := exec.Command(bin, "--snapshot", path, "--version", "1.22.0")
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit when --tool flag is missing")
	}
}

func TestMain_MissingVersionFlagExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	snap := makeSnap()
	path := writeJSON(t, snap)

	cmd := exec.Command(bin, "--snapshot", path, "--tool", "go")
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit when --version flag is missing")
	}
}

func TestMain_PinWritesAnnotation(t *testing.T) {
	bin := buildBinary(t)
	snap := makeSnap()
	path := writeJSON(t, snap)

	cmd := exec.Command(bin, "--snapshot", path, "--tool", "go", "--version", "1.22.0")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("expected success pinning tool: %v\n%s", err, out)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read updated snapshot: %v", err)
	}
	var updated snapshot.Snapshot
	if err := json.Unmarshal(data, &updated); err != nil {
		t.Fatalf("failed to unmarshal updated snapshot: %v", err)
	}
	if updated.Annotations == nil {
		t.Fatal("expected annotations to be set after pin")
	}
}

func TestMain_UnknownToolExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	snap := makeSnap()
	path := writeJSON(t, snap)

	cmd := exec.Command(bin, "--snapshot", path, "--tool", "nonexistent", "--version", "1.0.0")
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit for unknown tool")
	}
}
