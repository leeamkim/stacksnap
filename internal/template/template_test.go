package template_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stacksnap/internal/snapshot"
	"github.com/stacksnap/internal/template"
)

func makeTemplate() template.Template {
	return template.Template{
		Name: "test-template",
		Vars: []template.TemplateVar{
			{Name: "ENV", Description: "deployment environment", Default: "development"},
			{Name: "REGION", Description: "cloud region", Default: "us-east-1"},
		},
		Snapshot: snapshot.Snapshot{
			Env: map[string]string{
				"APP_ENV":    "{{ENV}}",
				"AWS_REGION": "{{REGION}}",
			},
			Tools: []snapshot.Tool{
				{Name: "go", Version: "1.22.0"},
			},
		},
	}
}

func TestSave_WritesJSONFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tmpl.json")
	tmpl := makeTemplate()

	if err := template.Save(tmpl, template.SaveOptions{Path: path}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var got template.Template
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Name != tmpl.Name {
		t.Errorf("name: got %q, want %q", got.Name, tmpl.Name)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tmpl.json")
	tmpl := makeTemplate()

	_ = template.Save(tmpl, template.SaveOptions{Path: path})
	loaded, err := template.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Vars) != len(tmpl.Vars) {
		t.Errorf("vars count: got %d, want %d", len(loaded.Vars), len(tmpl.Vars))
	}
}

func TestLoad_MissingFileReturnsError(t *testing.T) {
	_, err := template.Load("/no/such/file.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestInstantiate_ReplacesVars(t *testing.T) {
	tmpl := makeTemplate()
	snap := template.Instantiate(tmpl, map[string]string{"ENV": "production"})

	if snap.Env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", snap.Env["APP_ENV"], "production")
	}
	// REGION not provided — should fall back to default
	if snap.Env["AWS_REGION"] != "us-east-1" {
		t.Errorf("AWS_REGION: got %q, want %q", snap.Env["AWS_REGION"], "us-east-1")
	}
}

func TestInstantiate_UsesDefaultWhenNoValues(t *testing.T) {
	tmpl := makeTemplate()
	snap := template.Instantiate(tmpl, nil)

	if snap.Env["APP_ENV"] != "development" {
		t.Errorf("APP_ENV: got %q, want %q", snap.Env["APP_ENV"], "development")
	}
}

func TestInstantiate_PreservesTools(t *testing.T) {
	tmpl := makeTemplate()
	snap := template.Instantiate(tmpl, nil)

	if len(snap.Tools) != len(tmpl.Snapshot.Tools) {
		t.Fatalf("tools count: got %d, want %d", len(snap.Tools), len(tmpl.Snapshot.Tools))
	}
	if snap.Tools[0].Name != "go" || snap.Tools[0].Version != "1.22.0" {
		t.Errorf("tool: got %+v, want {Name:go Version:1.22.0}", snap.Tools[0])
	}
}
