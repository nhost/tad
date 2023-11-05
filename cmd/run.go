package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/nhost/tad/runbook"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

//go:embed basic.tpl
var tplBasic string

func LoadTOML(filepath string) (*runbook.Runbook, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	dec := toml.NewDecoder(f).DisallowUnknownFields()
	var runbook runbook.Runbook
	if err := dec.Decode(&runbook); err != nil {
		return nil, fmt.Errorf("failed to decode TOML: %w", err)
	}

	return &runbook, nil
}

func RunCommand() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:      "run",
		Usage:     "Runs a runbook",
		ArgsUsage: "<runbook>",
		Action: func(cCtx *cli.Context) error {
			if cCtx.Args().Len() != 1 {
				return fmt.Errorf("you need to specify a runbook") //nolint:goerr113
			}

			rb, err := LoadTOML(cCtx.Args().First())
			if err != nil {
				return err
			}

			interpreter, err := runbook.NewNodeREPL()
			if err != nil {
				return fmt.Errorf("failed to create interpreter: %w", err)
			}

			if err := interpreter.Start(cCtx.Context); err != nil {
				return fmt.Errorf("failed to start interpreter: %w", err)
			}

			if err := rb.Run(cCtx.Context, interpreter); err != nil {
				_ = rb.Render(tplBasic, os.Stdout, interpreter.Syntax())
				return fmt.Errorf("failed to run runbook: %w", err)
			}

			if err := rb.Render(tplBasic, os.Stdout, interpreter.Syntax()); err != nil {
				return fmt.Errorf("failed to render runbook: %w", err)
			}

			return nil
		},
	}
}
