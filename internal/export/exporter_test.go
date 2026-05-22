package export_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stacksnap/internal/export"
	"github.com/stacksnap/internal/snapshot"
)

func makeTestSnapshot() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		CapturedAt: time.Now(),
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
		},
		Env: map[string]string{
			"GOPATH": "/home/user/go",
		},
	}
}

func TestExport_WritesJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	snap := makeTestSnapshot()

	opts := export.Options{
		OutputDir: tmpDir,
		Format:    export.FormatJSON,
		Filename:  "test-snap.json",
	}

	outPath, err := export.Export(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatalf("expected output file to exist at %s", outPath)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	var result snapshot.Snapshot
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(result.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(result.Tools))
	}
}

func TestExport_DefaultFilenameContainsTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	snap := makeTestSnapshot()

	opts := export.DefaultOptions()
	opts.OutputDir = tmpDir

	outPath, err := export.Export(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	base := filepath.Base(outPath)
	if !strings.HasPrefix(base, "stacksnap-") {
		t.Errorf("expected filename to start with 'stacksnap-', got %s", base)
	}
	if !strings.HasSuffix(base, ".json") {
		t.Errorf("expected filename to end with '.json', got %s", base)
	}
}

func TestExport_UnsupportedFormatReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	snap := makeTestSnapshot()

	opts := export.Options{
		OutputDir: tmpDir,
		Format:    export.Format("xml"),
		Filename:  "snap.xml",
	}

	_, err := export.Export(snap, opts)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}
