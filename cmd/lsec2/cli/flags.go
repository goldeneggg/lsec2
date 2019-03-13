package cli

import (
	"os"

	"github.com/urfave/cli"
)

var (
	// command line options for getting ec2 instances
	ec2Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "region",
			Usage: "assign target region",
		},
		cli.StringFlag{
			Name:  "state, s",
			Usage: "filter only assigned state",
		},
	}

	// command line options for printing informations
	printFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "show-build",
			Usage: "show build info",
		},
		cli.BoolFlag{
			Name:  "print-header, H",
			Usage: "print list header",
		},
		cli.BoolFlag{
			Name:  "private-ip, p",
			Usage: "print only private ip",
		},
		cli.BoolFlag{
			Name:  "with-color, c",
			Usage: "print state with color",
		},
		cli.StringFlag{
			Name:  "delimiter, d",
			Usage: "column delemeter for printed informations",
			Value: "\t",
		},
		cli.StringFlag{
			Name:  "coldef",
			Usage: "path of coldef.yml",
			Value: os.Getenv("HOME") + string(os.PathSeparator) + ".lsec2" + string(os.PathSeparator) + "coldef.yml",
		},
	}
)
