package main_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stacksnap/internal/snapshot"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "stacksnap-tag")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func writeJSON(t *testing.T, dir string, s *snapshot.Snapshot) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	data, _ := json.MarshalIndent(s, "", "  ")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func makeSnap(tags ...string) *snapshot.Snapshot {
	return &snapshot.Snapshot{Tags: tags}
}

func TestMain_BuildsBinary(t *testing.T) {
	buildBinary(t)
}

func TestMain_MissingSnapshotFlagExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--list")
	if err := cmd.Run(); err == nil {
		t.Error("expected non-zero exit when --snapshot is missing")
	}
}

func TestMain_ListPrintsTags(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	path := writeJSON(t, dir, makeSnap("go", "docker"))
	out, err := exec.Command(bin, "--snapshot", path, "--list").Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(string(out))
	if !strings.Contains(got, "docker") || !strings.Contains(got, "go") {
		t.Errorf("expected tags in output, got: %q", got)
	}
}

func TestMain_AddTagUpdatesFile(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	path := writeJSON(t, dir, makeSnap())
	if out, err := exec.Command(bin, "--snapshot", path, "--add", "backend").CombinedOutput(); err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	data, _ := os.ReadFile(path)
	var s snapshot.Snapshot
	_ = json.Unmarshal(data, &s)
	if len(s.Tags) != 1 || s.Tags[0] != "backend" {
		t.Errorf("expected [backend], got %v", s.Tags)
	}
}

func TestMain_InvalidTagExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	path := writeJSON(t, dir, makeSnap())
	cmd := exec.Command(bin, "--snapshot", path, "--add", "bad tag!")
	if err := cmd.Run(); err == nil {
		t.Error("expected non-zero exit for invalid tag name")
	}
}
