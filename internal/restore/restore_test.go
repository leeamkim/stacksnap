package restore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

func writeSnapshotFile(t *testing.T, snap *snapshot.Snapshot) string {
	t.Helper()
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	path := filepath.Join(t.TempDir(), "snapshot.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func makeSnap(tools []snapshot.Tool) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		CapturedAt: time.Now(),
		Tools:      tools,
		Env:        map[string]string{},
	}
}

func TestLoadSnapshot_Valid(t *testing.T) {
	snap := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.22.0"}})
	path := writeSnapshotFile(t, snap)

	got, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Tools) != 1 || got.Tools[0].Name != "go" {
		t.Errorf("unexpected tools: %+v", got.Tools)
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_DryRunDoesNotError(t *testing.T) {
	snap := makeSnap([]snapshot.Tool{
		{Name: "go", Version: "1.22.0"},
		{Name: "node", Version: "20.11.0"},
	})
	if err := Apply(snap, Options{DryRun: true, Verbose: true}); err != nil {
		t.Fatalf("dry-run Apply returned error: %v", err)
	}
}

func TestInstallCommand_Unsupported(t *testing.T) {
	_, err := installCommand("unknowntool", "1.0")
	if err == nil {
		t.Fatal("expected error for unsupported tool")
	}
}

func TestInstallCommand_KnownTools(t *testing.T) {
	tools := []string{"go", "node", "python", "ruby"}
	for _, name := range tools {
		cmd, err := installCommand(name, "1.0.0")
		if err != nil {
			t.Errorf("%s: unexpected error: %v", name, err)
		}
		if cmd == "" {
			t.Errorf("%s: expected non-empty command", name)
		}
	}
}
