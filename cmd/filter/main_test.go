package main_test

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
	bin := filepath.Join(t.TempDir(), "filter")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func writeJSON(t *testing.T, snap snapshot.Snapshot) string {
	t.Helper()
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func makeSnap() snapshot.Snapshot {
	return snapshot.Snapshot{
		CapturedAt: time.Now(),
		OS:         "darwin",
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.0.0"},
		},
		Env: map[string]string{"HOME": "/home/user"},
	}
}

func TestMain_BuildsBinary(t *testing.T) {
	buildBinary(t)
}

func TestMain_MissingSnapshotFlagExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit when --snapshot is missing")
	}
}

func TestMain_NoEnvFlagExcludesEnv(t *testing.T) {
	bin := buildBinary(t)
	snap := makeSnap()
	snapshotFile := writeJSON(t, snap)

	cmd := exec.Command(bin, "--snapshot", snapshotFile, "--no-env")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	var result snapshot.Snapshot
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}
	if len(result.Env) != 0 {
		t.Errorf("expected no env vars, got %v", result.Env)
	}
}

func TestMain_TagsFilterTools(t *testing.T) {
	bin := buildBinary(t)
	snap := makeSnap()
	snapshotFile := writeJSON(t, snap)

	cmd := exec.Command(bin, "--snapshot", snapshotFile, "--tags", "go")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	var result snapshot.Snapshot
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}
	if len(result.Tools) != 1 || result.Tools[0].Name != "go" {
		t.Errorf("expected only 'go' tool, got %v", result.Tools)
	}
}
