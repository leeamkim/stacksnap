package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/stacksnap/internal/snapshot"
)

// Format defines the output format for the snapshot export.
type Format string

const (
	FormatJSON Format = "json"
	FormatTOML Format = "toml"
)

// Options holds configuration for the export operation.
type Options struct {
	OutputDir string
	Format    Format
	Filename  string
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		OutputDir: ".",
		Format:    FormatJSON,
		Filename:  "",
	}
}

// Export writes a snapshot to disk using the given options.
func Export(snap *snapshot.Snapshot, opts Options) (string, error) {
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	filename := opts.Filename
	if filename == "" {
		timestamp := time.Now().Format("20060102-150405")
		filename = fmt.Sprintf("stacksnap-%s.%s", timestamp, string(opts.Format))
	}

	outPath := filepath.Join(opts.OutputDir, filename)

	var data []byte
	var err error

	switch opts.Format {
	case FormatJSON:
		data, err = marshalJSON(snap)
	default:
		return "", fmt.Errorf("unsupported format: %s", opts.Format)
	}

	if err != nil {
		return "", fmt.Errorf("marshalling snapshot: %w", err)
	}

	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return "", fmt.Errorf("writing snapshot file: %w", err)
	}

	return outPath, nil
}

func marshalJSON(snap *snapshot.Snapshot) ([]byte, error) {
	return json.MarshalIndent(snap, "", "  ")
}
