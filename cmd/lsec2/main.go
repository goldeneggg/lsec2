package main

import (
	"fmt"
	"os"

	"github.com/goldeneggg/lsec2/cmd/lsec2/cli"
)

const (
	exitStsOk = iota
	exitStsNg
)

var (
	sts int
)

func main() {
	defer finalize()

	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		sts = exitStsNg
	}
}

func finalize() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Panic: %v\n", err)
		sts = exitStsNg
	}

	os.Exit(sts)
}
