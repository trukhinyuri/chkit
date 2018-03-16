package namespace

import (
	"fmt"

	"github.com/containerum/chkit/cmd/util"
	"gopkg.in/urfave/cli.v2"
)

// GetNamespace -- commmand 'get' entity data
var GetNamespace = &cli.Command{
	Name:        "ns",
	Description: `show namespace or namespace list`,
	Usage:       `Shows namespace data or namespace list`,
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		if ctx.NArg() > 0 {
			name := ctx.Args().First()
			ns, err := client.GetNamespace(name)
			if err != nil {
				return err
			}
			fmt.Println(ns.RenderTable())
		}
		return nil
	},
}