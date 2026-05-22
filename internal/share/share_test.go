package share_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/stacksnap/internal/share"
	"github.com/user/stacksnap/internal/snapshot"
)

func makeSnap() snapshot.Snapshot {
	return snapshot.Snapshot{
		CapturedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Tools: []snapshot.Tool{
			{Name: "go", Version: "1.22.0"},
			{Name: "node", Version: "20.11.0"},
		},
		Env: map[string]string{
			"GOPATH": "/home/user/go",
		},
	}
}

func TestEncode_HasSchemePrefix(t *testing.T) {
	snap := makeSnap()
	result, err := share.Encode(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(result, share.SchemePrefix) {
		t.Errorf("expected prefix %q, got %q", share.SchemePrefix, result[:min(len(result), 20)])
	}
}

func TestEncode_Decode_RoundTrip(t *testing.T) {
	original := makeSnap()
	encoded, err := share.Encode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := share.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded.Tools) != len(original.Tools) {
		t.Fatalf("expected %d tools, got %d", len(original.Tools), len(decoded.Tools))
	}
	for i, tool := range decoded.Tools {
		if tool.Name != original.Tools[i].Name || tool.Version != original.Tools[i].Version {
			t.Errorf("tool[%d] mismatch: got %+v, want %+v", i, tool, original.Tools[i])
		}
	}
}

func TestDecode_WithoutPrefix(t *testing.T) {
	original := makeSnap()
	encoded, _ := share.Encode(original)
	stripped := strings.TrimPrefix(encoded, share.SchemePrefix)

	decoded, err := share.Decode(stripped)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(decoded.Tools) != len(original.Tools) {
		t.Errorf("expected %d tools, got %d", len(original.Tools), len(decoded.Tools))
	}
}

func TestDecode_InvalidBase64(t *testing.T) {
	_, err := share.Decode("stacksnap://!!!notbase64!!!")
	if err == nil {
		t.Error("expected error for invalid base64, got nil")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
