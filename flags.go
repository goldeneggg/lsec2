package main

import "github.com/urfave/cli"

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "region",
		Usage: "assign target region",
	},
}

var showFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "print-header, H",
		Usage: "print list header",
	},
	cli.BoolFlag{
		Name:  "private-ip, p",
		Usage: "print only private ip",
	},
	cli.BoolFlag{
		Name:  "show-build",
		Usage: "show build info",
	},
	cli.BoolFlag{
		Name:  "with-color, c",
		Usage: "print state with color",
	},
}
