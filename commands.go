package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

import "github.com/goldeneggg/lsec2/awsec2"

func run(args []string) error {
	app := cli.NewApp()

	app.Name = "lsec2"
	app.Author = "goldeneggg"
	app.Version = VERSION
	app.Usage = "Listing information of aws ec2 instances"
	//app.Flags = globalFlags
	app.Flags = append(globalFlags, showFlags...)
	//app.Commands = commands
	app.Action = action

	return app.Run(args)
}

func action(c *cli.Context) {
	if c.IsSet("show-build") {
		showBuildInfo(c)
		return
	}

	client := &awsec2.Client{
		Region: c.String("region"),
		Tags:   c.Args(),
	}

	if c.IsSet("H") {
		client.PrintHeader = c.Bool("H")
	}

	if c.IsSet("p") {
		client.OnlyPrivateIP = c.Bool("p")
		client.PrintHeader = false
	}

	if c.IsSet("c") {
		client.WithColor = c.Bool("c")
	}

	if c.IsSet("s") {
		client.StateName = c.String("s")
	}

	if err := client.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		sts = exitStsNg
	}
}

func showBuildInfo(c *cli.Context) {
	fmt.Printf("build-date: %v\n", buildDate)
	fmt.Printf("build-commit: %v\n", buildCommit)
}
