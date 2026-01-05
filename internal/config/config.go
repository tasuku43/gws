package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Root string
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gws", "config.yaml"), nil
}

func Load(path string) (Config, error) {
	if path == "" {
		var err error
		path, err = DefaultPath()
		if err != nil {
			return Config{}, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}

	root, err := parseRoot(data)
	if err != nil {
		return Config{}, err
	}

	return Config{Root: root}, nil
}

func parseRoot(data []byte) (string, error) {
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.HasPrefix(trimmed, "root:") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "root:"))
		if value == "" {
			return "", nil
		}
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) >= 2 {
			value = strings.TrimSuffix(strings.TrimPrefix(value, "\""), "\"")
		}
		if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") && len(value) >= 2 {
			value = strings.TrimSuffix(strings.TrimPrefix(value, "'"), "'")
		}
		return value, nil
	}

	return "", nil
}
