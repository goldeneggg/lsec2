package cli

import (
	"fmt"
	"os"

	"github.com/goldeneggg/lsec2/awsec2"
	"github.com/goldeneggg/lsec2/cmd/lsec2/version"
	"github.com/goldeneggg/lsec2/printer"
	"github.com/urfave/cli"
)

var (
	// BuildDate have a built date
	BuildDate string
	// BuildCommit have a latest commit hash
	BuildCommit string
	// GoVersion have a go version number used for build
	GoVersion string
)

// Run run a cli command
func Run(args []string) error {
	app := cli.NewApp()

	app.Name = "lsec2"
	app.Author = "goldeneggg"
	app.Version = version.VERSION
	app.Usage = "Listing information of aws ec2 instances"
	app.Flags = append(ec2Flags, printFlags...)
	app.Action = action

	return app.Run(args)
}

func action(c *cli.Context) error {
	if c.IsSet("show-build") {
		showBuildInfo(c)
		return nil
	}

	if err := newPrinter(c).PrintAll(newClient(c)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	return nil
}

func newClient(c *cli.Context) *awsec2.Client {
	client := awsec2.NewClient(c.String("region"), c.String("profile"))

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
		printer.Delimiter = c.String("d")
	}

	return printer
}

func showBuildInfo(c *cli.Context) {
	fmt.Printf("build-date: %v\n", BuildDate)
	fmt.Printf("build-commit: %v\n", BuildCommit)
	fmt.Printf("go-version: %v\n", GoVersion)
}
