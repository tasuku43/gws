package gitcmd

import (
	"context"
	"fmt"
	"strings"
)

// ForEachRefContains returns ref names that contain the given commit.
func ForEachRefContains(ctx context.Context, dir, commit, refPrefix string) ([]string, error) {
	args := []string{"for-each-ref", "--format=%(refname)", "--contains", commit}
	if strings.TrimSpace(refPrefix) != "" {
		args = append(args, refPrefix)
	}
	res, err := Run(ctx, args, Options{Dir: dir})
	if err != nil {
		if strings.TrimSpace(res.Stderr) != "" {
			return nil, fmt.Errorf("git for-each-ref failed: %w: %s", err, strings.TrimSpace(res.Stderr))
		}
		return nil, fmt.Errorf("git for-each-ref failed: %w", err)
	}
	output := strings.TrimSpace(res.Stdout)
	if output == "" {
		return nil, nil
	}
	lines := strings.Split(output, "\n")
	var refs []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		refs = append(refs, line)
	}
	return refs, nil
}
