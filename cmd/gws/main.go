package main

import (
	"fmt"
	"os"

	"github.com/tasuku43/gws/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
