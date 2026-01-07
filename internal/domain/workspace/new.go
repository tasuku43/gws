package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tasuku43/gws/internal/core/gitcmd"
)

func New(ctx context.Context, rootDir string, workspaceID string) (string, error) {
	if err := validateWorkspaceID(ctx, workspaceID); err != nil {
		return "", err
	}
	if rootDir == "" {
		return "", fmt.Errorf("root directory is required")
	}

	wsDir := filepath.Join(rootDir, "workspaces", workspaceID)
	if exists, err := pathExists(wsDir); err != nil {
		return "", err
	} else if exists {
		return "", fmt.Errorf("workspace already exists: %s", wsDir)
	}

	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		return "", fmt.Errorf("create workspace dir: %w", err)
	}

	return wsDir, nil
}

func validateWorkspaceID(ctx context.Context, workspaceID string) error {
	if workspaceID == "" {
		return fmt.Errorf("workspace id is required")
	}
	_, err := gitcmd.Run(ctx, []string{"check-ref-format", "--branch", workspaceID}, gitcmd.Options{})
	if err != nil {
		return fmt.Errorf("invalid workspace id: %w", err)
	}
	return nil
}

func pathExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("path is not a directory: %s", path)
	}
	return true, nil
}
