package restore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/stacksnap/internal/snapshot"
)

// Options configures the restore behaviour.
type Options struct {
	DryRun  bool
	Verbose bool
}

// DefaultOptions returns sensible defaults for restore.
func DefaultOptions() Options {
	return Options{
		DryRun:  false,
		Verbose: false,
	}
}

// LoadSnapshot reads and deserialises a snapshot from the given file path.
func LoadSnapshot(path string) (*snapshot.Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("restore: reading snapshot file: %w", err)
	}
	var snap snapshot.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("restore: parsing snapshot JSON: %w", err)
	}
	return &snap, nil
}

// Apply attempts to install or configure each tool recorded in the snapshot.
// When opts.DryRun is true the commands are printed but not executed.
func Apply(snap *snapshot.Snapshot, opts Options) error {
	for _, tool := range snap.Tools {
		cmd, err := installCommand(tool.Name, tool.Version)
		if err != nil {
			if opts.Verbose {
				fmt.Printf("[skip] %s: %v\n", tool.Name, err)
			}
			continue
		}
		if opts.Verbose || opts.DryRun {
			fmt.Printf("[run] %s\n", cmd)
		}
		if opts.DryRun {
			continue
		}
		if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
			return fmt.Errorf("restore: installing %s: %w\noutput: %s", tool.Name, err, out)
		}
	}
	return nil
}

// installCommand returns a shell command string that installs the given tool
// at the requested version. Returns an error when the tool is unsupported.
func installCommand(name, version string) (string, error) {
	switch name {
	case "go":
		return fmt.Sprintf("go install golang.org/dl/go%s@latest && go%s download", version, version), nil
	case "node":
		return fmt.Sprintf("fnm install %s", version), nil
	case "python":
		return fmt.Sprintf("pyenv install -s %s", version), nil
	case "ruby":
		return fmt.Sprintf("rbenv install -s %s", version), nil
	default:
		return "", fmt.Errorf("no install strategy for %q", name)
	}
}
