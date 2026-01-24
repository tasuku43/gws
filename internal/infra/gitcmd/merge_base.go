package gitcmd

import (
	"context"
	"fmt"
	"strings"
)

func IsAncestor(ctx context.Context, dir, ancestor, descendant string) (bool, error) {
	ancestor = strings.TrimSpace(ancestor)
	descendant = strings.TrimSpace(descendant)
	if ancestor == "" {
		return false, fmt.Errorf("ancestor is required")
	}
	if descendant == "" {
		return false, fmt.Errorf("descendant is required")
	}

	res, err := Run(ctx, []string{"merge-base", "--is-ancestor", ancestor, descendant}, Options{Dir: dir})
	if err == nil {
		return true, nil
	}
	if res.ExitCode == 1 {
		return false, nil
	}
	if strings.TrimSpace(res.Stderr) != "" {
		return false, fmt.Errorf("git merge-base --is-ancestor failed: %w: %s", err, strings.TrimSpace(res.Stderr))
	}
	return false, fmt.Errorf("git merge-base --is-ancestor failed: %w", err)
}
