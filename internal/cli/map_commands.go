package cli

import (
	"github.com/leotech-dev/leoflow/internal/common"
	"github.com/leotech-dev/leoflow/internal/mapper"
	"github.com/leotech-dev/leoflow/internal/mapper/jq"
	nfmapper "github.com/numaproj/numaflow-go/pkg/mapper"
	"github.com/urfave/cli/v2"
)

var mapCmd = &cli.Command{
	Name:      "map",
	UsageText: "leoflow map command [argumets...]",
	Subcommands: []*cli.Command{
		mapJqCommand,
	},
}

var mapJqCommand = &cli.Command{
	Name:      "jq",
	UsageText: "leoflow map jq",
	Usage:     "Run jq expressions on the input JSON data",
	Action:    mapHandler(&jq.JqMapper{}),
}

func mapHandler(m mapper.Mapper) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if i, ok := m.(common.Initializable); ok {
			i.Init()
		}

		return nfmapper.NewServer(m).Start(ctx.Context)
	}
}
