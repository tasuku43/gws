package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

type Entry struct {
	WorkspaceID   string
	WorkspacePath string
}

func List(rootDir string) ([]Entry, []error, error) {
	wsRoot := filepath.Join(rootDir, "workspaces")
	info, err := os.Stat(wsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("workspaces path is not a directory: %s", wsRoot)
	}

	entries, err := os.ReadDir(wsRoot)
	if err != nil {
		return nil, nil, err
	}

	var results []Entry

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		wsID := entry.Name()
		wsPath := filepath.Join(wsRoot, wsID)

		result := Entry{
			WorkspaceID:   wsID,
			WorkspacePath: wsPath,
		}
		results = append(results, result)
	}

	return results, nil, nil
}
