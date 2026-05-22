package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain_BuildsBinary verifies the cmd compiles without errors.
func TestMain_BuildsBinary(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "stacksnap-snapshot")

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, string(out))
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Fatal("expected binary to exist after build")
	}
}

// TestMain_RunsAndProducesFile verifies the binary runs and writes a snapshot.
func TestMain_RunsAndProducesFile(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "stacksnap-snapshot")

	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, string(out))
	}

	cmd := exec.Command(binPath, "--out", tmpDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, string(out))
	}

	output := string(out)
	if !strings.Contains(output, "Snapshot written to:") {
		t.Errorf("expected success message in output, got: %s", output)
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read output dir: %v", err)
	}

	found := false
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected a .json snapshot file in output directory")
	}
}

// TestMain_UnsupportedFormatExitsNonZero verifies bad format flag causes failure.
func TestMain_UnsupportedFormatExitsNonZero(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "stacksnap-snapshot")

	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, string(out))
	}

	cmd := exec.Command(binPath, "--format", "yaml", "--out", tmpDir)
	err := cmd.Run()
	if err == nil {
		t.Error("expected non-zero exit for unsupported format, got nil")
	}
}
