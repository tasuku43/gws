package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFile(t *testing.T) {
	rootDir := t.TempDir()
	if _, err := Load(rootDir); err == nil {
		t.Fatalf("expected error for missing templates file")
	}
}

func TestLoadAndNames(t *testing.T) {
	rootDir := t.TempDir()
	path := filepath.Join(rootDir, FileName)
	data := []byte(`
templates:
  app:
    repos:
      - git@github.com:org/app.git
  zzz:
    repos:
      - git@github.com:org/zzz.git
  legacy:
    repos:
      - repo: git@github.com:org/legacy.git
`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write templates: %v", err)
	}
	file, err := Load(rootDir)
	if err != nil {
		t.Fatalf("load templates: %v", err)
	}
	names := Names(file)
	want := []string{"app", "legacy", "zzz"}
	if len(names) != len(want) {
		t.Fatalf("expected %d names, got %d", len(want), len(names))
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("name mismatch at %d: got %q want %q", i, names[i], want[i])
		}
	}
	legacy, ok := file.Templates["legacy"]
	if !ok || len(legacy.Repos) != 1 {
		t.Fatalf("legacy template not loaded")
	}
}
