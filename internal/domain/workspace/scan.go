package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tasuku43/gws/internal/core/gitcmd"
	"github.com/tasuku43/gws/internal/domain/repospec"
)

func ScanRepos(ctx context.Context, wsDir string) ([]Repo, []error, error) {
	entries, err := os.ReadDir(wsDir)
	if err != nil {
		return nil, nil, err
	}
	var repos []Repo
	var warnings []error
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == ".gws" {
			continue
		}
		repoPath := filepath.Join(wsDir, entry.Name())
		repo, warn, ok := inspectRepo(ctx, repoPath, entry.Name())
		if !ok {
			if warn != nil {
				warnings = append(warnings, warn)
			}
			continue
		}
		if warn != nil {
			warnings = append(warnings, warn)
		}
		repos = append(repos, repo)
	}
	return repos, warnings, nil
}

func inspectRepo(ctx context.Context, repoPath, alias string) (Repo, error, bool) {
	gitDir, err := gitRevParse(ctx, repoPath, "--git-dir")
	if err != nil {
		return Repo{}, fmt.Errorf("skip %s: not a git repo", repoPath), false
	}
	commonDir, err := gitRevParse(ctx, repoPath, "--git-common-dir")
	if err != nil {
		return Repo{}, fmt.Errorf("skip %s: %v", repoPath, err), false
	}

	absGitDir := resolveGitPath(repoPath, gitDir)
	absCommonDir := resolveGitPath(repoPath, commonDir)

	storePath := ""
	if absGitDir != absCommonDir {
		storePath = absCommonDir
	}

	branch := readBranch(ctx, repoPath)
	repoSpec, repoKey, warn := readRepoSpec(ctx, repoPath)

	repo := Repo{
		Alias:        alias,
		RepoSpec:     repoSpec,
		RepoKey:      repoKey,
		StorePath:    storePath,
		WorktreePath: repoPath,
		Branch:       branch,
	}
	return repo, warn, true
}

func gitRevParse(ctx context.Context, repoPath, arg string) (string, error) {
	res, err := gitcmd.Run(ctx, []string{"rev-parse", arg}, gitcmd.Options{Dir: repoPath})
	if err != nil {
		if strings.TrimSpace(res.Stderr) != "" {
			return "", fmt.Errorf("git rev-parse %s failed: %w: %s", arg, err, strings.TrimSpace(res.Stderr))
		}
		return "", fmt.Errorf("git rev-parse %s failed: %w", arg, err)
	}
	return strings.TrimSpace(res.Stdout), nil
}

func resolveGitPath(repoPath, value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Clean(filepath.Join(repoPath, value))
}

func readBranch(ctx context.Context, repoPath string) string {
	res, err := gitcmd.Run(ctx, []string{"symbolic-ref", "--short", "HEAD"}, gitcmd.Options{Dir: repoPath})
	if err == nil {
		return strings.TrimSpace(res.Stdout)
	}
	return ""
}

func readRepoSpec(ctx context.Context, repoPath string) (string, string, error) {
	res, err := gitcmd.Run(ctx, []string{"remote", "get-url", "origin"}, gitcmd.Options{Dir: repoPath})
	if err != nil {
		if strings.TrimSpace(res.Stderr) != "" {
			return "", "", fmt.Errorf("origin remote missing: %s", strings.TrimSpace(res.Stderr))
		}
		return "", "", fmt.Errorf("origin remote missing")
	}
	remoteURL := strings.TrimSpace(res.Stdout)
	if remoteURL == "" {
		return "", "", fmt.Errorf("origin remote is empty")
	}
	spec, err := repospec.Normalize(remoteURL)
	if err != nil {
		return remoteURL, "", fmt.Errorf("origin remote invalid: %s", err)
	}
	return remoteURL, spec.RepoKey, nil
}
