package main

import (
	"fmt"
	"os"
)

import (
	"github.com/goldeneggg/lsec2/constants"
)

var sts int

func main() {
	defer finalize()

	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		sts = constants.ExitStsNg
	}
}

func finalize() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Panic: %v\n", err)
		sts = constants.ExitStsNg
	}

	os.Exit(sts)
}
