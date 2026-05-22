package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stacksnap/internal/history"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap(version string) snapshot.Snapshot {
	return snapshot.Snapshot{
		Timestamp: time.Now().UTC(),
		Tools: []snapshot.Tool{
			{Name: "go", Version: version},
		},
		Env: map[string]string{"GOPATH": "/home/user/go"},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	snap := makeSnap("1.22.0")

	entry, err := history.Save(snap, "initial", dir)
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	if entry.Label != "initial" {
		t.Errorf("expected label 'initial', got %q", entry.Label)
	}
	if entry.ID == "" {
		t.Error("expected non-empty ID")
	}

	path := filepath.Join(dir, entry.ID+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", path)
	}
}

func TestSave_DefaultLabelIsID(t *testing.T) {
	dir := t.TempDir()
	entry, err := history.Save(makeSnap("1.21.0"), "", dir)
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if entry.Label != entry.ID {
		t.Errorf("expected label == ID when label is empty, got label=%q id=%q", entry.Label, entry.ID)
	}
}

func TestList_ReturnsSortedEntries(t *testing.T) {
	dir := t.TempDir()

	for _, v := range []string{"1.20.0", "1.21.0", "1.22.0"} {
		time.Sleep(10 * time.Millisecond) // ensure distinct timestamps
		_, err := history.Save(makeSnap(v), v, dir)
		if err != nil {
			t.Fatalf("Save(%s) error: %v", v, err)
		}
	}

	entries, err := history.List(dir)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Timestamp.Before(entries[i-1].Timestamp) {
			t.Errorf("entries not sorted: index %d before index %d", i, i-1)
		}
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	entries, err := history.List(dir)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
