package cli

import (
	"github.com/urfave/cli/v2"
)

func NewApplication() *cli.App {
	return &cli.App{
		Name:      "Leoflow",
		Usage:     "a set of Numaflow plugins by Leotech",
		UsageText: "leoflow [global options] command [command options] [arguments...]",
		Commands: []*cli.Command{
			mapCmd,
		},
	}
}
