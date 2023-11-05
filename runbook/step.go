package runbook

import (
	"context"
	"errors"
	"fmt"
)

type Step struct {
	Pre         string `toml:"pre"`
	Post        string `toml:"post"`
	Exec        string `toml:"exec"`
	IgnoreError bool   `toml:"ignore_error"`
	Stdout      string `toml:"_"`
	Stderr      string `toml:"_"`
	Err         error  `toml:"_"`
}

func (s *Step) Run(
	_ context.Context,
	interpreter Interpreter,
) error {
	stdout, stderr, err := interpreter.Exec(s.Exec)
	s.Stdout = stdout
	s.Stderr = stderr

	runtimeError := &RuntimeError{} //nolint:exhaustruct
	if err != nil && !(s.IgnoreError && errors.As(err, &runtimeError)) {
		s.Err = err
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}
