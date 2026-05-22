package snapshot

import (
	"os"
	"strings"
	"time"
)

// Tool represents a detected CLI tool and its version.
type Tool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Snapshot holds a point-in-time capture of the local dev environment.
type Snapshot struct {
	CapturedAt time.Time         `json:"captured_at"`
	Tools      []Tool            `json:"tools"`
	Env        map[string]string `json:"env"`
}

// CaptureOptions controls what is included in the snapshot.
type CaptureOptions struct {
	// EnvPrefixes filters environment variables to those matching these prefixes.
	// If empty, no env vars are captured.
	EnvPrefixes []string
}

// Capture collects the current dev stack state and returns a Snapshot.
func Capture(opts CaptureOptions) (*Snapshot, error) {
	tools, err := detectTools()
	if err != nil {
		return nil, err
	}

	env := filterEnv(os.Environ(), opts.EnvPrefixes)

	return &Snapshot{
		CapturedAt: time.Now(),
		Tools:      tools,
		Env:        env,
	}, nil
}

// filterEnv returns env vars whose keys match any of the given prefixes.
func filterEnv(environ []string, prefixes []string) map[string]string {
	result := make(map[string]string)
	if len(prefixes) == 0 {
		return result
	}
	for _, entry := range environ {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) {
				result[key] = val
				break
			}
		}
	}
	return result
}
