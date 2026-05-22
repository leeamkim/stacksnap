package template

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/stacksnap/internal/snapshot"
)

// TemplateVar represents a placeholder in a template snapshot.
type TemplateVar struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

// Template wraps a snapshot with variable definitions for reuse.
type Template struct {
	Name      string                 `json:"name"`
	Vars      []TemplateVar          `json:"vars,omitempty"`
	Snapshot  snapshot.Snapshot      `json:"snapshot"`
}

// SaveOptions controls how a template is written to disk.
type SaveOptions struct {
	Path string
}

// DefaultSaveOptions returns sensible defaults for saving a template.
func DefaultSaveOptions() SaveOptions {
	return SaveOptions{
		Path: "stacksnap-template.json",
	}
}

// Save writes a Template to disk as JSON.
func Save(t Template, opts SaveOptions) error {
	if opts.Path == "" {
		opts = DefaultSaveOptions()
	}
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("template: marshal: %w", err)
	}
	if err := os.WriteFile(opts.Path, b, 0o644); err != nil {
		return fmt.Errorf("template: write %s: %w", opts.Path, err)
	}
	return nil
}

// Load reads a Template from a JSON file.
func Load(path string) (Template, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Template{}, fmt.Errorf("template: read %s: %w", path, err)
	}
	var t Template
	if err := json.Unmarshal(b, &t); err != nil {
		return Template{}, fmt.Errorf("template: unmarshal: %w", err)
	}
	return t, nil
}

// Instantiate replaces {{VAR}} placeholders in env values using the provided
// values map, falling back to each TemplateVar's Default.
func Instantiate(t Template, values map[string]string) snapshot.Snapshot {
	snap := t.Snapshot
	resolved := make(map[string]string, len(snap.Env))
	for k, v := range snap.Env {
		resolved[k] = expandVars(v, t.Vars, values)
	}
	snap.Env = resolved
	return snap
}

func expandVars(s string, vars []TemplateVar, values map[string]string) string {
	for _, v := range vars {
		placeholder := "{{" + v.Name + "}}"
		val, ok := values[v.Name]
		if !ok || val == "" {
			val = v.Default
		}
		s = strings.ReplaceAll(s, placeholder, val)
	}
	return s
}
