package filter_test

import (
	"testing"
	"time"

	"github.com/stacksnap/internal/filter"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() snapshot.Snapshot {
	return snapshot.Snapshot{
		CapturedAt: time.Now(),
		OS:         "linux",
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
			{Name: "docker", Version: "24.0.5"},
		},
		Env: map[string]string{
			"GOPATH": "/home/user/go",
			"HOME":   "/home/user",
		},
	}
}

func TestApply_IncludesAllByDefault(t *testing.T) {
	snap := makeSnap()
	out := filter.Apply(snap, filter.DefaultOptions())

	if out.OS != snap.OS {
		t.Errorf("expected OS %q, got %q", snap.OS, out.OS)
	}
	if len(out.Tools) != len(snap.Tools) {
		t.Errorf("expected %d tools, got %d", len(snap.Tools), len(out.Tools))
	}
	if len(out.Env) != len(snap.Env) {
		t.Errorf("expected %d env vars, got %d", len(snap.Env), len(out.Env))
	}
}

func TestApply_ExcludesEnv(t *testing.T) {
	snap := makeSnap()
	opts := filter.DefaultOptions()
	opts.Env = false
	out := filter.Apply(snap, opts)

	if len(out.Env) != 0 {
		t.Errorf("expected no env vars, got %d", len(out.Env))
	}
	if len(out.Tools) != len(snap.Tools) {
		t.Error("tools should still be present")
	}
}

func TestApply_ExcludesTools(t *testing.T) {
	snap := makeSnap()
	opts := filter.DefaultOptions()
	opts.Tools = false
	out := filter.Apply(snap, opts)

	if len(out.Tools) != 0 {
		t.Errorf("expected no tools, got %d", len(out.Tools))
	}
}

func TestApply_FilterByTag(t *testing.T) {
	snap := makeSnap()
	opts := filter.DefaultOptions()
	opts.Tags = []string{"go", "docker"}
	out := filter.Apply(snap, opts)

	if len(out.Tools) != 2 {
		t.Errorf("expected 2 tools after tag filter, got %d", len(out.Tools))
	}
}

func TestApply_TagMatchIsCaseInsensitive(t *testing.T) {
	snap := makeSnap()
	opts := filter.DefaultOptions()
	opts.Tags = []string{"GO"}
	out := filter.Apply(snap, opts)

	if len(out.Tools) != 1 || out.Tools[0].Name != "go" {
		t.Errorf("expected 1 tool matching 'GO', got %v", out.Tools)
	}
}

func TestApply_ExcludesOS(t *testing.T) {
	snap := makeSnap()
	opts := filter.DefaultOptions()
	opts.OS = false
	out := filter.Apply(snap, opts)

	if out.OS != "" {
		t.Errorf("expected empty OS, got %q", out.OS)
	}
}
