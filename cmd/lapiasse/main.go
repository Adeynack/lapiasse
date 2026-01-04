package main

import (
	"fmt"
	"os"
)

func main() {
	// If LAPIASSE_DEFAULT_CMD is set and no subcommand is provided, use the command from the env var.
	// e.g.: LAPIASSE_DEFAULT_CMD=serve lapiasse is equivalent to lapiasse serve
	if cmd := os.Getenv("LAPIASSE_DEFAULT_CMD"); cmd != "" && len(os.Args) == 1 {
		os.Args = append(os.Args, cmd)
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
