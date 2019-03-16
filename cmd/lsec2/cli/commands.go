package cli

import (
	"fmt"

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

	pr := newPrinter(c)
	client, err := newClient(c)
	if err != nil {
		return err
	}

	if err := pr.PrintAll(client); err != nil {
		return err
	}

	return nil
}

func newPrinter(c *cli.Context) *printer.Printer {
	return printer.NewPrinter(c.String("d"), c.Bool("H"), c.Bool("p"), c.Bool("c"), nil)
}

func newClient(c *cli.Context) (*awsec2.Client, error) {
	return awsec2.NewClient(c.String("region"), c.String("s"), c.String("profile"), c.Args())
}

func showBuildInfo(c *cli.Context) {
	fmt.Printf("build-date: %v\n", BuildDate)
	fmt.Printf("build-commit: %v\n", BuildCommit)
	fmt.Printf("go-version: %v\n", GoVersion)
}
