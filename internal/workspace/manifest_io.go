package workspace

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func LoadManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func WriteManifest(path string, manifest Manifest) error {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

func TouchLastUsed(manifest *Manifest, now time.Time) {
	manifest.LastUsedAt = now.UTC().Format(time.RFC3339)
}

func AddRepo(manifest *Manifest, repo Repo) bool {
	for _, existing := range manifest.Repos {
		if existing.Alias == repo.Alias {
			return false
		}
		if existing.RepoKey != "" && existing.RepoKey == repo.RepoKey {
			return false
		}
	}

	manifest.Repos = append(manifest.Repos, repo)
	return true
}
