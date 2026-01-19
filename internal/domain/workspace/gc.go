package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/tasuku43/gwst/internal/core/gitcmd"
)

type GCRule struct {
	Name   string
	Reason string
	Match  func(info GCRepoInfo) (bool, error)
}

type GCRepoInfo struct {
	Status             RepoStatus
	Head               string
	OriginDefaultRef   string
	OriginContainsRefs []string
}

type GCRepoCandidate struct {
	Alias   string
	RepoKey string
	Branch  string
	Reasons []string
}

type GCCandidate struct {
	WorkspaceID   string
	WorkspacePath string
	Description   string
	Repos         []GCRepoCandidate
}

type GCResult struct {
	Scanned    int
	Candidates []GCCandidate
	Skipped    int
	Warnings   []error
}

func DefaultGCRules() []GCRule {
	return []GCRule{
		{
			Name:   "pushed-to-origin",
			Reason: "already pushed to origin",
			Match: func(info GCRepoInfo) (bool, error) {
				return len(info.OriginContainsRefs) > 0, nil
			},
		},
		{
			Name:   "merged-into-origin-default",
			Reason: "merged into origin default branch",
			Match: func(info GCRepoInfo) (bool, error) {
				if strings.TrimSpace(info.OriginDefaultRef) == "" {
					return false, nil
				}
				return refListContains(info.OriginContainsRefs, info.OriginDefaultRef), nil
			},
		},
	}
}

func ScanGC(ctx context.Context, rootDir string, rules []GCRule) (GCResult, error) {
	entries, warnings, err := List(rootDir)
	if err != nil {
		return GCResult{}, err
	}
	if len(rules) == 0 {
		rules = DefaultGCRules()
	}
	result := GCResult{
		Scanned:  len(entries),
		Warnings: append([]error(nil), warnings...),
	}
	for _, entry := range entries {
		candidate, warnings, ok := evaluateWorkspaceForGC(ctx, rootDir, entry, rules)
		if len(warnings) > 0 {
			result.Warnings = append(result.Warnings, warnings...)
		}
		if ok {
			result.Candidates = append(result.Candidates, candidate)
		} else {
			result.Skipped++
		}
	}
	return result, nil
}

func evaluateWorkspaceForGC(ctx context.Context, rootDir string, entry Entry, rules []GCRule) (GCCandidate, []error, bool) {
	status, err := Status(ctx, rootDir, entry.WorkspaceID)
	if err != nil {
		return GCCandidate{}, []error{fmt.Errorf("workspace %s status: %w", entry.WorkspaceID, err)}, false
	}
	var warnings []error
	if len(status.Warnings) > 0 {
		for _, warning := range status.Warnings {
			warnings = append(warnings, fmt.Errorf("workspace %s: %w", entry.WorkspaceID, warning))
		}
	}

	candidate := GCCandidate{
		WorkspaceID:   entry.WorkspaceID,
		WorkspacePath: entry.WorkspacePath,
		Description:   entry.Description,
	}

	for _, repoStatus := range status.Repos {
		if repoStatus.Error != nil {
			warnings = append(warnings, fmt.Errorf("workspace %s repo %s: %v", entry.WorkspaceID, repoStatus.Alias, repoStatus.Error))
		}
		if repoExcludedForGC(repoStatus) {
			return GCCandidate{}, warnings, false
		}
		info, err := buildGCRepoInfo(ctx, repoStatus)
		if err != nil {
			warnings = append(warnings, fmt.Errorf("workspace %s repo %s: %v", entry.WorkspaceID, repoStatus.Alias, err))
			return GCCandidate{}, warnings, false
		}
		reasons, err := matchGCRules(info, rules)
		if err != nil {
			warnings = append(warnings, fmt.Errorf("workspace %s repo %s: %v", entry.WorkspaceID, repoStatus.Alias, err))
			return GCCandidate{}, warnings, false
		}
		if len(reasons) == 0 {
			return GCCandidate{}, warnings, false
		}
		candidate.Repos = append(candidate.Repos, GCRepoCandidate{
			Alias:   repoStatus.Alias,
			RepoKey: repoStatus.RepoKey,
			Branch:  repoStatus.Branch,
			Reasons: reasons,
		})
	}

	if len(candidate.Repos) == 0 {
		return GCCandidate{}, warnings, false
	}
	return candidate, warnings, true
}

func repoExcludedForGC(status RepoStatus) bool {
	if status.Error != nil {
		return true
	}
	if status.Dirty {
		return true
	}
	if status.Detached || status.HeadMissing {
		return true
	}
	if status.AheadCount > 0 {
		return true
	}
	return false
}

func buildGCRepoInfo(ctx context.Context, status RepoStatus) (GCRepoInfo, error) {
	if strings.TrimSpace(status.WorktreePath) == "" {
		return GCRepoInfo{}, fmt.Errorf("worktree path is required")
	}
	head, err := gitcmd.RevParse(ctx, status.WorktreePath, "HEAD")
	if err != nil {
		return GCRepoInfo{}, err
	}
	originDefaultRef, ok, err := gitcmd.SymbolicRef(ctx, status.WorktreePath, "refs/remotes/origin/HEAD")
	if err != nil {
		return GCRepoInfo{}, err
	}
	if !ok {
		originDefaultRef = ""
	}
	containsRefs, err := gitcmd.ForEachRefContains(ctx, status.WorktreePath, head, "refs/remotes/origin")
	if err != nil {
		return GCRepoInfo{}, err
	}
	return GCRepoInfo{
		Status:             status,
		Head:               head,
		OriginDefaultRef:   originDefaultRef,
		OriginContainsRefs: containsRefs,
	}, nil
}

func matchGCRules(info GCRepoInfo, rules []GCRule) ([]string, error) {
	var reasons []string
	for _, rule := range rules {
		if rule.Match == nil {
			continue
		}
		matched, err := rule.Match(info)
		if err != nil {
			return nil, err
		}
		if matched {
			reasons = append(reasons, rule.Reason)
		}
	}
	return reasons, nil
}

func refListContains(refs []string, target string) bool {
	for _, ref := range refs {
		if ref == target {
			return true
		}
	}
	return false
}
