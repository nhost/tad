package main

import (
	"log"
	"os"

	"github.com/nhost/tad/cmd"
	"github.com/urfave/cli/v2"
)

var Version string

func main() {
	app := &cli.App{
		Version: Version,
		Commands: []*cli.Command{
			cmd.RunCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
