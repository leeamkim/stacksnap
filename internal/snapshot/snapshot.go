package snapshot

import (
	"os"
	"runtime"
	"time"
)

// Snapshot represents a captured state of the local dev stack.
type Snapshot struct {
	CapturedAt  time.Time         `json:"captured_at"`
	Hostname    string            `json:"hostname"`
	OS          string            `json:"os"`
	Arch        string            `json:"arch"`
	Environment map[string]string `json:"environment"`
	Tools       []ToolInfo        `json:"tools"`
}

// ToolInfo holds version and path details for a detected tool.
type ToolInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Path    string `json:"path"`
}

// Capture collects the current dev stack state and returns a Snapshot.
func Capture(envKeys []string) (*Snapshot, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	env := make(map[string]string)
	for _, key := range envKeys {
		if val, ok := os.LookupEnv(key); ok {
			env[key] = val
		}
	}

	tools, err := detectTools()
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		CapturedAt:  time.Now().UTC(),
		Hostname:    hostname,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		Environment: env,
		Tools:       tools,
	}, nil
}
