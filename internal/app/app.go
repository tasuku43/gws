package app

import (
	"flag"
	"os"

	"github.com/tasuku43/gws/internal/paths"
)

// Run is a placeholder for the CLI entrypoint.
func Run() error {
	fs := flag.NewFlagSet("gws", flag.ContinueOnError)
	var rootFlag string
	fs.StringVar(&rootFlag, "root", "", "override gws root")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	_, err := paths.ResolveRoot(rootFlag)
	return err
}
