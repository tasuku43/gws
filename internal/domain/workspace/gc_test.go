package workspace_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tasuku43/gwst/internal/domain/repo"
	"github.com/tasuku43/gwst/internal/domain/workspace"
)

func TestScanGCCandidates(t *testing.T) {
	t.Setenv("GIT_AUTHOR_NAME", "gwst")
	t.Setenv("GIT_AUTHOR_EMAIL", "gwst@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "gwst")
	t.Setenv("GIT_COMMITTER_EMAIL", "gwst@example.com")

	ctx := context.Background()
	tmp := t.TempDir()
	rootDir := filepath.Join(tmp, "gwst")

	remoteBase := filepath.Join(tmp, "remotes")
	remotePath := filepath.Join(remoteBase, "org", "repo.git")
	if err := os.MkdirAll(filepath.Dir(remotePath), 0o755); err != nil {
		t.Fatalf("mkdir remote: %v", err)
	}
	runGit(t, "", "init", "--bare", remotePath)

	seedDir := filepath.Join(tmp, "seed")
	runGit(t, "", "init", seedDir)
	runGit(t, seedDir, "checkout", "-b", "main")
	if err := os.WriteFile(filepath.Join(seedDir, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write seed file: %v", err)
	}
	runGit(t, seedDir, "add", ".")
	runGit(t, seedDir, "commit", "-m", "init")
	runGit(t, seedDir, "remote", "add", "origin", remotePath)
	runGit(t, seedDir, "push", "origin", "main")
	runGit(t, "", "--git-dir", remotePath, "symbolic-ref", "HEAD", "refs/heads/main")

	configPath := filepath.Join(tmp, "gitconfig")
	fileURL := "file://" + filepath.ToSlash(remoteBase) + "/"
	configData := fmt.Sprintf("[url \"%s\"]\n\tinsteadOf = https://example.com/\n", fileURL)
	if err := os.WriteFile(configPath, []byte(configData), 0o644); err != nil {
		t.Fatalf("write gitconfig: %v", err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", configPath)
	t.Setenv("GIT_CONFIG_SYSTEM", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_TERMINAL_PROMPT", "0")

	repoSpec := "https://example.com/org/repo.git"
	if _, err := repo.Get(ctx, rootDir, repoSpec); err != nil {
		t.Fatalf("repo get: %v", err)
	}

	if _, err := workspace.New(ctx, rootDir, "WS-1"); err != nil {
		t.Fatalf("workspace new: %v", err)
	}
	if _, err := workspace.Add(ctx, rootDir, "WS-1", repoSpec, "", true); err != nil {
		t.Fatalf("workspace add: %v", err)
	}

	if _, err := workspace.New(ctx, rootDir, "WS-2"); err != nil {
		t.Fatalf("workspace new: %v", err)
	}
	if _, err := workspace.Add(ctx, rootDir, "WS-2", repoSpec, "", true); err != nil {
		t.Fatalf("workspace add: %v", err)
	}
	worktreePath := workspace.WorktreePath(rootDir, "WS-2", "repo")
	if err := os.WriteFile(filepath.Join(worktreePath, "change.txt"), []byte("change\n"), 0o644); err != nil {
		t.Fatalf("write worktree file: %v", err)
	}
	runGit(t, worktreePath, "add", ".")
	runGit(t, worktreePath, "commit", "-m", "local change")

	result, err := workspace.ScanGC(ctx, rootDir, nil)
	if err != nil {
		t.Fatalf("scan gc: %v", err)
	}
	if result.Scanned != 2 {
		t.Fatalf("unexpected scanned count: got %d want %d", result.Scanned, 2)
	}
	if len(result.Candidates) != 1 {
		t.Fatalf("unexpected candidates: got %d want %d (warnings: %v)", len(result.Candidates), 1, result.Warnings)
	}
	if result.Candidates[0].WorkspaceID != "WS-1" {
		t.Fatalf("unexpected candidate id: got %s want %s", result.Candidates[0].WorkspaceID, "WS-1")
	}
	if len(result.Candidates[0].Repos) == 0 {
		t.Fatalf("candidate repos missing")
	}
	if !containsString(result.Candidates[0].Repos[0].Reasons, "already pushed to origin") {
		t.Fatalf("expected reason to include pushed-to-origin")
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
