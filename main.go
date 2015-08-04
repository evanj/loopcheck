package main

import (
	"fmt"
	"os"

	"github.com/evanj/loopcheck/libloopcheck"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "loopcheck (files to check)")
		fmt.Fprintln(os.Stderr, "  checks a single package for loop variables that might escape")
		os.Exit(1)
	}

	if err := libloopcheck.CheckFiles(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to check files: %s\n", err)
		os.Exit(1)
	}
}
