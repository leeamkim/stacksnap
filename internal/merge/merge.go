package merge

import (
	"fmt"
	"sort"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

// Strategy defines how conflicting tool versions are resolved.
type Strategy string

const (
	StrategyBase    Strategy = "base"     // prefer base snapshot
	StrategyIncoming Strategy = "incoming" // prefer incoming snapshot
	StrategyNewest  Strategy = "newest"   // prefer the newer version string
)

// Result holds the merged snapshot and a log of any conflicts resolved.
type Result struct {
	Snapshot  snapshot.Snapshot
	Conflicts []string
}

// Merge combines two snapshots into one according to the given strategy.
// Tools present in only one snapshot are always included. Conflicts arise
// when the same tool exists in both snapshots with different versions.
func Merge(base, incoming snapshot.Snapshot, strategy Strategy) (Result, error) {
	if strategy != StrategyBase && strategy != StrategyIncoming && strategy != StrategyNewest {
		return Result{}, fmt.Errorf("merge: unknown strategy %q", strategy)
	}

	toolIndex := make(map[string]snapshot.Tool)
	var conflicts []string

	for _, t := range base.Tools {
		toolIndex[t.Name] = t
	}

	for _, inc := range incoming.Tools {
		existing, exists := toolIndex[inc.Name]
		if !exists {
			toolIndex[inc.Name] = inc
			continue
		}
		if existing.Version == inc.Version {
			continue
		}
		// Conflict: same tool, different version
		conflicts = append(conflicts, fmt.Sprintf("%s: base=%s incoming=%s", inc.Name, existing.Version, inc.Version))
		switch strategy {
		case StrategyIncoming:
			toolIndex[inc.Name] = inc
		case StrategyNewest:
			if inc.Version > existing.Version {
				toolIndex[inc.Name] = inc
			}
		// StrategyBase: keep existing — no-op
		}
	}

	// Collect and sort tools by name for deterministic output.
	merged := make([]snapshot.Tool, 0, len(toolIndex))
	for _, t := range toolIndex {
		merged = append(merged, t)
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Name < merged[j].Name
	})

	// Merge env: base wins on conflict, incoming adds new keys
	env := make(map[string]string, len(base.Env))
	for k, v := range base.Env {
		env[k] = v
	}
	for k, v := range incoming.Env {
		if _, exists := env[k]; !exists {
			env[k] = v
		}
	}

	return Result{
		Snapshot: snapshot.Snapshot{
			Timestamp: time.Now().UTC(),
			Tools:     merged,
			Env:       env,
		},
		Conflicts: conflicts,
	}, nil
}
