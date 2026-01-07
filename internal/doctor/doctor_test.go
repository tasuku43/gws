package doctor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tasuku43/gws/internal/workspace"
)

func TestCheckFindsIssues(t *testing.T) {
	rootDir := t.TempDir()
	now := time.Now().UTC()

	wsDir := filepath.Join(rootDir, "workspaces", "WS1")
	manifestDir := filepath.Join(wsDir, ".gws")
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		t.Fatalf("mkdir ws: %v", err)
	}
	manifestPath := filepath.Join(manifestDir, "manifest.yaml")
	manifest := workspace.Manifest{
		SchemaVersion: 1,
		WorkspaceID:   "WS1",
		Repos: []workspace.Repo{
			{
				Alias:        "app",
				WorktreePath: filepath.Join(wsDir, "missing"),
			},
		},
	}
	if err := workspace.WriteManifest(manifestPath, manifest); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	wsLock := filepath.Join(manifestDir, "lock")
	if err := os.WriteFile(wsLock, []byte("lock"), 0o644); err != nil {
		t.Fatalf("write ws lock: %v", err)
	}
	if err := os.Chtimes(wsLock, now.Add(-48*time.Hour), now.Add(-48*time.Hour)); err != nil {
		t.Fatalf("chtimes ws lock: %v", err)
	}

	repoOK := filepath.Join(rootDir, "bare", "example.com", "org", "ok.git")
	if err := os.MkdirAll(repoOK, 0o755); err != nil {
		t.Fatalf("mkdir repo ok: %v", err)
	}
	configOK := []byte("[remote \"origin\"]\n\turl = git@github.com:org/ok.git\n")
	if err := os.WriteFile(filepath.Join(repoOK, "config"), configOK, 0o644); err != nil {
		t.Fatalf("write config ok: %v", err)
	}
	repoLock := filepath.Join(repoOK, ".gws", "lock")
	if err := os.MkdirAll(filepath.Dir(repoLock), 0o755); err != nil {
		t.Fatalf("mkdir repo lock dir: %v", err)
	}
	if err := os.WriteFile(repoLock, []byte("lock"), 0o644); err != nil {
		t.Fatalf("write repo lock: %v", err)
	}
	if err := os.Chtimes(repoLock, now.Add(-48*time.Hour), now.Add(-48*time.Hour)); err != nil {
		t.Fatalf("chtimes repo lock: %v", err)
	}

	repoNoRemote := filepath.Join(rootDir, "bare", "example.com", "org", "noremote.git")
	if err := os.MkdirAll(repoNoRemote, 0o755); err != nil {
		t.Fatalf("mkdir repo noremote: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoNoRemote, "config"), []byte("[core]\n\trepositoryformatversion = 0\n"), 0o644); err != nil {
		t.Fatalf("write config noremote: %v", err)
	}

	result, err := Check(context.Background(), rootDir, now)
	if err != nil {
		t.Fatalf("doctor check: %v", err)
	}
	kinds := map[string]int{}
	for _, issue := range result.Issues {
		kinds[issue.Kind]++
	}
	if kinds["stale_lock"] == 0 {
		t.Fatalf("expected stale_lock issue")
	}
	if kinds["missing_worktree"] == 0 {
		t.Fatalf("expected missing_worktree issue")
	}
	if kinds["missing_remote"] == 0 {
		t.Fatalf("expected missing_remote issue")
	}
}

func TestCheckRootLayout(t *testing.T) {
	rootDir := t.TempDir()
	now := time.Now().UTC()

	result, err := Check(context.Background(), rootDir, now)
	if err != nil {
		t.Fatalf("doctor check: %v", err)
	}
	kinds := map[string]int{}
	for _, issue := range result.Issues {
		kinds[issue.Kind]++
	}
	if kinds["missing_root_dir"] == 0 {
		t.Fatalf("expected missing_root_dir issues")
	}
	if kinds["missing_root_file"] == 0 {
		t.Fatalf("expected missing_root_file issues")
	}
}
