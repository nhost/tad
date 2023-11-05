package runbook

import (
	"context"
	"fmt"
	"io"
	"text/template"
)

type Interpreter interface {
	Start(ctx context.Context) error
	Close() error
	Exec(cmd string) (string, string, error)
}

type Global struct {
	Interpreter string `toml:"interpreter"`
	Pre         string `toml:"pre"`
}

type Runbook struct {
	Global *Global `toml:"global"`
	Steps  []*Step `toml:"steps"`
}

func New(global *Global, steps []*Step) *Runbook {
	return &Runbook{
		Global: global,
		Steps:  steps,
	}
}

func (r *Runbook) Run(ctx context.Context, interpreter Interpreter) error {
	for i, s := range r.Steps {
		if err := s.Run(ctx, interpreter); err != nil {
			// fmt.Println("err: ", err)
			return fmt.Errorf("failed to run step %d: %w", i+1, err)
		}
	}

	return nil
}

func (r *Runbook) Render(tpl string, wr io.Writer, syntax string) error {
	t, err := template.New("runbook").Parse(tpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	if err := t.Execute(
		wr,
		map[string]any{
			"Global": r.Global,
			"Steps":  r.Steps,
			"Syntax": syntax,
		},
	); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
