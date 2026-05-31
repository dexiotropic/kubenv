package main

import (
	"fmt"
	"os"

	"github.com/dexiotropic/kubenv/internal/cli"
)

func main() {
	if err := cli.RunKubectlPlugin(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
