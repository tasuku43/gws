package paths

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tasuku43/gws/internal/config"
)

const defaultRootDir = "gws"

func ResolveRoot(flagRoot string) (string, error) {
	if flagRoot != "" {
		return normalizeRoot(flagRoot)
	}

	envRoot := os.Getenv("GWS_ROOT")
	if envRoot != "" {
		return normalizeRoot(envRoot)
	}

	cfg, err := config.Load("")
	if err != nil {
		return "", err
	}
	if cfg.Root != "" {
		return normalizeRoot(cfg.Root)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultRootDir), nil
}

func normalizeRoot(path string) (string, error) {
	expanded, err := expandHome(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(expanded), nil
}

func expandHome(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
	}

	return path, nil
}
