package main

import (
	"fmt"
	"os"

	"github.com/nhost/tad/cmd"
	"github.com/urfave/cli/v2"
)

var Version string

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name:    "tad",
		Version: Version,
		Commands: []*cli.Command{
			cmd.RunCommand(),
			cmd.RunMD(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
