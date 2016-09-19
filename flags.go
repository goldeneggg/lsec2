package main

import "github.com/urfave/cli"

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "region",
		Value: "ap-northeast-1",
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
}
