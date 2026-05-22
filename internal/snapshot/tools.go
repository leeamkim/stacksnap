package snapshot

import (
	"os/exec"
	"strings"
)

// knownTools lists common dev tools to probe.
var knownTools = []string{
	"go",
	"node",
	"npm",
	"python3",
	"docker",
	"git",
	"make",
}

// detectTools probes known tools and returns their info.
func detectTools() ([]ToolInfo, error) {
	var tools []ToolInfo

	for _, name := range knownTools {
		path, err := exec.LookPath(name)
		if err != nil {
			// tool not found; skip
			continue
		}

		version := getVersion(name)
		tools = append(tools, ToolInfo{
			Name:    name,
			Version: version,
			Path:    path,
		})
	}

	return tools, nil
}

// versionArgs maps tool names to their version flag.
var versionArgs = map[string][]string{
	"go":      {"version"},
	"node":    {"--version"},
	"npm":     {"--version"},
	"python3": {"--version"},
	"docker":  {"--version"},
	"git":     {"--version"},
	"make":    {"--version"},
}

func getVersion(name string) string {
	args, ok := versionArgs[name]
	if !ok {
		return "unknown"
	}

	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "unknown"
	}

	line := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
	return line
}
