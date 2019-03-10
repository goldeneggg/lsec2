package main

import (
	"fmt"
	"os"

	"github.com/goldeneggg/lsec2/awsec2"
	"github.com/goldeneggg/lsec2/printer"
	"github.com/urfave/cli"
)

func run(args []string) error {
	app := cli.NewApp()

	app.Name = "lsec2"
	app.Author = "goldeneggg"
	app.Version = VERSION
	app.Usage = "Listing information of aws ec2 instances"
	app.Flags = append(ec2Flags, printFlags...)
	app.Action = action

	return app.Run(args)
}

func action(c *cli.Context) {
	if c.IsSet("show-build") {
		showBuildInfo(c)
		return
	}

	if err := newPrinter(c).PrintAll(newClient(c)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		sts = exitStsNg
	}
}

func newClient(c *cli.Context) *awsec2.Client {
	client := awsec2.NewClient(c.String("region"), nil)

	client.Tags = c.Args()
	if c.IsSet("s") {
		client.StateName = c.String("s")
	}

	return client
}

func newPrinter(c *cli.Context) *printer.Printer {
	printer := printer.NewPrinter(nil)

	if c.IsSet("H") {
		printer.PrintHeader = c.Bool("H")
	}
	if c.IsSet("p") {
		printer.OnlyPrivateIP = c.Bool("p")
		printer.PrintHeader = false
	}
	if c.IsSet("c") {
		printer.WithColor = c.Bool("c")
	}
	if c.IsSet("d") {
		printer.Delimeter = c.String("d")
	}
	if c.IsSet("coldef") {
		printer.ColDef = c.String("coldef")
	}

	return printer
}

func showBuildInfo(c *cli.Context) {
	fmt.Printf("build-date: %v\n", buildDate)
	fmt.Printf("build-commit: %v\n", buildCommit)
}
