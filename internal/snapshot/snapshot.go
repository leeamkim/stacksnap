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

// EnvAllowlist defines which environment variable prefixes are captured.
var EnvAllowlist = []string{
	"GOPATH", "GOROOT", "GOVERSION",
	"NODE_", "NPM_", "NVM_",
	"JAVA_", "PYTHON",
	"PATH",
}

// Capture collects the current environment and detected tools into a Snapshot.
func Capture() (*Snapshot, error) {
	tools, err := detectTools()
	if err != nil {
		return nil, err
	}
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Tools:      tools,
		Env:        filterEnv(os.Environ()),
	}, nil
}

// filterEnv returns only the environment variables matching the allowlist.
func filterEnv(environ []string) map[string]string {
	result := make(map[string]string)
	for _, entry := range environ {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		if isAllowed(key) {
			result[key] = val
		}
	}
	return result
}

func isAllowed(key string) bool {
	for _, prefix := range EnvAllowlist {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
