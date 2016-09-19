package main

import (
	"github.com/goldeneggg/lsec2/awsec2"
	"github.com/urfave/cli"
)

var commands = []cli.Command{
	cli.Command{
		Name:   "show",
		Usage:  "show instance list",
		Flags:  showFlags,
		Action: action,
	},
}

func run(args []string) error {
	app := cli.NewApp()

	app.Name = "lsec2"
	app.Author = "goldeneggg"
	app.Version = VERSION
	app.Usage = "Listing information of aws ec2 instances"
	app.Flags = globalFlags
	app.Commands = commands

	return app.Run(args)
}

func action(c *cli.Context) {
	opt := &awsec2.Opt{
		Region: c.GlobalString("region"),
		Tags:   c.Args(),
	}

	if c.IsSet("H") {
		opt.PrintHeader = c.Bool("H")
	}

	if c.IsSet("p") {
		opt.OnlyPrivateIP = c.Bool("p")
		opt.PrintHeader = false
	}

	sts = opt.Show()
}
