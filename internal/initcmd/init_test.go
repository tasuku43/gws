package initcmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCreatesAndSkips(t *testing.T) {
	rootDir := t.TempDir()

	result, err := Run(rootDir)
	if err != nil {
		t.Fatalf("init run: %v", err)
	}
	if len(result.CreatedDirs) == 0 {
		t.Fatalf("expected created dirs")
	}
	if len(result.CreatedFiles) == 0 {
		t.Fatalf("expected created files")
	}

	for _, dir := range []string{"bare", "src", "ws"} {
		if _, err := os.Stat(filepath.Join(rootDir, dir)); err != nil {
			t.Fatalf("missing dir %s: %v", dir, err)
		}
	}
	for _, file := range []string{"settings.yaml", "templates.yaml"} {
		if _, err := os.Stat(filepath.Join(rootDir, file)); err != nil {
			t.Fatalf("missing file %s: %v", file, err)
		}
	}

	second, err := Run(rootDir)
	if err != nil {
		t.Fatalf("second init run: %v", err)
	}
	if len(second.CreatedDirs) != 0 || len(second.CreatedFiles) != 0 {
		t.Fatalf("expected no created items on second run")
	}
	if len(second.SkippedDirs) == 0 || len(second.SkippedFiles) == 0 {
		t.Fatalf("expected skipped items on second run")
	}
}
