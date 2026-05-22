package snapshot

import (
	"testing"
)

func TestCapture_ReturnsSnapshot(t *testing.T) {
	snap, err := Capture([]string{"PATH", "HOME"})
	if err != nil {
		t.Fatalf("Capture() returned unexpected error: %v", err)
	}

	if snap == nil {
		t.Fatal("Capture() returned nil snapshot")
	}

	if snap.OS == "" {
		t.Error("expected OS to be populated")
	}

	if snap.Arch == "" {
		t.Error("expected Arch to be populated")
	}

	if snap.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestCapture_EnvFiltering(t *testing.T) {
	snap, err := Capture([]string{"PATH"})
	if err != nil {
		t.Fatalf("Capture() error: %v", err)
	}

	if _, ok := snap.Environment["PATH"]; !ok {
		t.Error("expected PATH to be present in environment map")
	}
}

func TestDetectTools_ReturnsSlice(t *testing.T) {
	tools, err := detectTools()
	if err != nil {
		t.Fatalf("detectTools() error: %v", err)
	}

	// We can't assert specific tools are present (CI may differ),
	// but each entry must have a non-empty name and path.
	for _, tool := range tools {
		if tool.Name == "" {
			t.Error("tool name should not be empty")
		}
		if tool.Path == "" {
			t.Errorf("tool %q path should not be empty", tool.Name)
		}
	}
}

func TestGetVersion_UnknownTool(t *testing.T) {
	v := getVersion("nonexistent-tool-xyz")
	if v != "unknown" {
		t.Errorf("expected 'unknown', got %q", v)
	}
}
