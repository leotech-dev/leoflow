package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/leotech-dev/leoflow/internal/cli"
)

func main() {
	err := cli.NewApplication().RunContext(context.Background(), os.Args)
	if err != nil {
		slog.Error("application error", "error", err)
	}
}
