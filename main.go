package main

import (
	"fmt"
	"os"
)

const (
	exitStsOk = iota
	exitStsNg
)

var (
	sts int

	version     string
	buildDate   string
	buildCommit string
)

func main() {
	defer finalize()

	if err := run(os.Args); err != nil {
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
