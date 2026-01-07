package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tasuku43/gws/internal/core/gitcmd"
)

func Remove(ctx context.Context, rootDir, workspaceID string) error {
	if workspaceID == "" {
		return fmt.Errorf("workspace id is required")
	}
	if rootDir == "" {
		return fmt.Errorf("root directory is required")
	}

	wsDir := filepath.Join(rootDir, "workspaces", workspaceID)
	if exists, err := pathExists(wsDir); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("workspace does not exist: %s", wsDir)
	}

	repos, warnings, err := ScanRepos(ctx, wsDir)
	if err != nil {
		return err
	}
	_ = warnings

	for _, repo := range repos {
		if repo.WorktreePath == "" {
			return fmt.Errorf("missing worktree path for alias %q", repo.Alias)
		}
		statusOut, statusErr := gitStatusPorcelain(ctx, repo.WorktreePath)
		if statusErr != nil {
			return fmt.Errorf("check status for %q: %w", repo.Alias, statusErr)
		}
		_, _, _, dirty, _, _, _, _, _, _ := parseStatusPorcelainV2(statusOut, "")
		if dirty {
			return fmt.Errorf("workspace has dirty changes: %s", repo.Alias)
		}
	}

	for _, repo := range repos {
		if repo.StorePath == "" {
			continue
		}
		if repo.WorktreePath == "" {
			return fmt.Errorf("missing worktree path for alias %q", repo.Alias)
		}
		gitcmd.Logf("git worktree remove %s", repo.WorktreePath)
		if _, err := gitcmd.Run(ctx, []string{"worktree", "remove", repo.WorktreePath}, gitcmd.Options{Dir: repo.StorePath}); err != nil {
			return fmt.Errorf("remove worktree %q: %w", repo.Alias, err)
		}
	}

	if err := os.RemoveAll(wsDir); err != nil {
		return fmt.Errorf("remove workspace dir: %w", err)
	}

	return nil
}
