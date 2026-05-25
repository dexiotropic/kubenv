package main

import (
	"fmt"
	"os"

	"github.com/dexiotropic/kubenv/internal/cmp"
)

func main() {
	if err := cmp.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
