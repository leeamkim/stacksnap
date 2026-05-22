package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "restore")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func writeJSON(t *testing.T, snap *snapshot.Snapshot) string {
	t.Helper()
	data, _ := json.Marshal(snap)
	p := filepath.Join(t.TempDir(), "snap.json")
	_ = os.WriteFile(p, data, 0o644)
	return p
}

func TestMain_BuildsBinary(t *testing.T) {
	buildBinary(t)
}

func TestMain_MissingSnapshotFlagExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit without --snapshot flag")
	}
	if !strings.Contains(string(out), "--snapshot") {
		t.Errorf("expected usage hint in output, got: %s", out)
	}
}

func TestMain_DryRunPrintsCommands(t *testing.T) {
	bin := buildBinary(t)
	snap := &snapshot.Snapshot{
		CapturedAt: time.Now(),
		Tools:      []snapshot.Tool{{Name: "go", Version: "1.22.0"}},
		Env:        map[string]string{},
	}
	path := writeJSON(t, snap)

	cmd := exec.Command(bin, "--snapshot", path, "--dry-run", "--verbose")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("dry-run exited non-zero: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "dry-run") {
		t.Errorf("expected dry-run output, got: %s", out)
	}
}
