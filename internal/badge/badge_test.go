package badge_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stacksnap/internal/badge"
	"github.com/stacksnap/internal/snapshot"
)

func makeSnap() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        "test-123",
		CreatedAt: time.Now(),
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
		},
		Env: map[string]string{"GOPATH": "/home/user/go"},
	}
}

func TestGenerate_ReturnsResult(t *testing.T) {
	snap := makeSnap()
	res, err := badge.Generate(snap, badge.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGenerate_URLContainsShieldsHost(t *testing.T) {
	res, err := badge.Generate(makeSnap(), badge.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(res.ShieldsURL, "https://img.shields.io/badge/") {
		t.Errorf("unexpected URL prefix: %s", res.ShieldsURL)
	}
}

func TestGenerate_MarkdownContainsURL(t *testing.T) {
	res, err := badge.Generate(makeSnap(), badge.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Markdown, res.ShieldsURL) {
		t.Errorf("markdown does not embed shields URL")
	}
}

func TestGenerate_NilSnapshotReturnsError(t *testing.T) {
	_, err := badge.Generate(nil, badge.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestGenerate_CustomLabelAppearsInURL(t *testing.T) {
	opts := badge.DefaultOptions()
	opts.Label = "myproject"
	res, err := badge.Generate(makeSnap(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.ShieldsURL, "myproject") {
		t.Errorf("expected custom label in URL, got: %s", res.ShieldsURL)
	}
}

func TestGenerate_StyleAppearsInURL(t *testing.T) {
	opts := badge.DefaultOptions()
	opts.Style = badge.StyleFlatSquare
	res, err := badge.Generate(makeSnap(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.ShieldsURL, "flat-square") {
		t.Errorf("expected style in URL, got: %s", res.ShieldsURL)
	}
}

func TestGenerate_GradeAndScorePopulated(t *testing.T) {
	res, err := badge.Generate(makeSnap(), badge.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Grade == "" {
		t.Error("expected non-empty grade")
	}
	if res.Score <= 0 {
		t.Errorf("expected positive score, got %d", res.Score)
	}
}
