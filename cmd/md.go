package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nhost/tad/runbook"
	"github.com/urfave/cli/v2"
)

func openMDFile(filepath string) (string, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(b), nil
}

func processStep(exec string) (string, bool) {
	lines := strings.Split(exec, "\n")
	ignoreError := strings.Contains(lines[0], "ignore_error")

	return strings.Join(lines[1:], "\n"), ignoreError
}

func getSteps(md string) []*runbook.Step {
	sections := strings.Split(md, "```")

	steps := make([]*runbook.Step, len(sections)/2) //nolint:gomnd

	i := 0
	for {
		j := i * 2 //nolint:gomnd

		pre := sections[j]
		exec, ignoreError := processStep(sections[j+1])

		var post string

		// if there is a third hanging we add it as post
		if j+3 == len(sections) {
			post = sections[j+2]
		}

		steps[i] = &runbook.Step{
			Pre:         pre,
			Post:        post,
			Exec:        exec,
			IgnoreError: ignoreError,
			Stdout:      "",
			Stderr:      "",
			Err:         nil,
		}

		if j+3 >= len(sections) {
			break
		}

		i++
	}

	return steps
}

func RunMD() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct
		Name:      "md",
		Usage:     "Runs a runbook in markdown format",
		ArgsUsage: "<markdown>",
		Action: func(cCtx *cli.Context) error {
			if cCtx.Args().Len() != 1 {
				return fmt.Errorf("you need to specify a runbook") //nolint:goerr113
			}

			md, err := openMDFile(cCtx.Args().First())
			if err != nil {
				return err
			}

			steps := getSteps(md)

			rb := runbook.New(
				&runbook.Global{
					Interpreter: "js",
					Pre:         "",
				},
				steps,
			)

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
