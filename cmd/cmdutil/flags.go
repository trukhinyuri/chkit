package cmdutil

import cli "gopkg.in/urfave/cli.v2"

var NamespaceFlag = &cli.StringFlag{
	Name:    "namespace",
	Aliases: []string{"n", "ns"},
	Usage:   "namespace to use",
}

var GetFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "file",
		Usage:   "file to write output",
		Aliases: []string{"f"},
	},
	&cli.StringFlag{
		Name:    "output",
		Usage:   "define output formats: yaml, json",
		Aliases: []string{"o"},
	},
	NamespaceFlag,
}

var DeleteFlags = []cli.Flag{
	NamespaceFlag,
}