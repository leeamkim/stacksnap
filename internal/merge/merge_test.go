package merge

import (
	"testing"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

func makeSnap(tools []snapshot.Tool, env map[string]string) snapshot.Snapshot {
	if env == nil {
		env = map[string]string{}
	}
	return snapshot.Snapshot{
		Timestamp: time.Now().UTC(),
		Tools:     tools,
		Env:       env,
	}
}

func TestMerge_NoConflicts(t *testing.T) {
	base := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.22.0"}}, nil)
	incoming := makeSnap([]snapshot.Tool{{Name: "node", Version: "20.0.0"}}, nil)

	res, err := Merge(base, incoming, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if len(res.Snapshot.Tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(res.Snapshot.Tools))
	}
}

func TestMerge_ConflictStrategyBase(t *testing.T) {
	base := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.21.0"}}, nil)
	incoming := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.22.0"}}, nil)

	res, err := Merge(base, incoming, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	for _, tool := range res.Snapshot.Tools {
		if tool.Name == "go" && tool.Version != "1.21.0" {
			t.Errorf("expected base version 1.21.0, got %s", tool.Version)
		}
	}
}

func TestMerge_ConflictStrategyIncoming(t *testing.T) {
	base := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.21.0"}}, nil)
	incoming := makeSnap([]snapshot.Tool{{Name: "go", Version: "1.22.0"}}, nil)

	res, err := Merge(base, incoming, StrategyIncoming)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, tool := range res.Snapshot.Tools {
		if tool.Name == "go" && tool.Version != "1.22.0" {
			t.Errorf("expected incoming version 1.22.0, got %s", tool.Version)
		}
	}
}

func TestMerge_EnvMerge(t *testing.T) {
	base := makeSnap(nil, map[string]string{"GO111MODULE": "on", "SHARED": "base"})
	incoming := makeSnap(nil, map[string]string{"NODE_ENV": "development", "SHARED": "incoming"})

	res, err := Merge(base, incoming, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Snapshot.Env["SHARED"] != "base" {
		t.Errorf("expected base to win on SHARED, got %s", res.Snapshot.Env["SHARED"])
	}
	if res.Snapshot.Env["NODE_ENV"] != "development" {
		t.Errorf("expected NODE_ENV from incoming, got %s", res.Snapshot.Env["NODE_ENV"])
	}
}

func TestMerge_UnknownStrategyReturnsError(t *testing.T) {
	base := makeSnap(nil, nil)
	incoming := makeSnap(nil, nil)

	_, err := Merge(base, incoming, Strategy("unknown"))
	if err == nil {
		t.Error("expected error for unknown strategy, got nil")
	}
}
