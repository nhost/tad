package runbook

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

//go:embed node_repl.js
var nodejs []byte

type NodeREPL struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
	writer io.WriteCloser
}

func NewNodeREPL() (*NodeREPL, error) {
	return &NodeREPL{
		cmd:    nil,
		stdout: nil,
		writer: nil,
		stderr: nil,
	}, nil
}

func (nr *NodeREPL) Syntax() string {
	return "js"
}

func (nr *NodeREPL) Start(ctx context.Context) error {
	f, err := os.CreateTemp(os.TempDir(), "node_repl.js")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if _, err := f.Write(nodejs); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	cmd := exec.CommandContext(ctx, "node", f.Name()) //nolint:gosec
	nr.cmd = cmd

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	nr.writer = stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	nr.stdout = stdout

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	nr.stderr = stderr

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}

	if _, _, err = nr.read(); err != nil {
		return fmt.Errorf("failed to read until initial prompt: %w", err)
	}

	return nil
}

func (nr *NodeREPL) Close() error {
	if _, _, err := nr.Exec("process.exit()"); err != nil {
		return fmt.Errorf("failed to exit repl: %w", err)
	}

	if err := nr.stdout.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}

	if err := nr.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	if err := nr.cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for cmd: %w", err)
	}

	return nil
}

func isUncaught(buf []byte) bool {
	return buf[0] == 'U' && buf[1] == 'n' && buf[2] == 'c' && buf[3] == 'a' &&
		buf[4] == 'u' &&
		buf[5] == 'g' &&
		buf[6] == 'h' &&
		buf[7] == 't'
}

func (nr *NodeREPL) read() (string, string, error) { //nolint:cyclop
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	buf := make([]byte, 1024) //nolint:gomnd
	for {
		n, err := nr.stdout.Read(buf)
		// fmt.Println(n, buf[:n], string(buf[:n]))
		if err != nil {
			return "", "", fmt.Errorf(
				"failed to read from stdout: %w. Last buffer data:\n %s",
				err,
				stdout.String(),
			)
		}

		// Uncaught
		if n > 10 && isUncaught(buf[:8]) {
			stderr.Write(buf[9:n])
			continue
		}

		// 0 is the null character we configured the repl to send back
		// make sure you don't mistake it with an empty int
		if n < 1024 && buf[n-1] == 0 {
			stdout.Write(buf[:n-1])
			break
		}

		if n == 4 && buf[0] == '.' && buf[1] == '.' && buf[2] == '.' {
			break
		}

		stdout.Write(buf[:n])
	}

	return stdout.String(), stderr.String(), nil
}

func (nr *NodeREPL) Exec(message string) (string, string, error) {
	msg := strings.ReplaceAll(message, "\n", "")
	_, err := nr.writer.Write([]byte(msg + "\n"))
	if err != nil {
		return "", "", fmt.Errorf("failed to write to stdin: %w", err)
	}

	o, e, err := nr.read()
	if err != nil {
		return "", "", err
	}

	if e != "" {
		return o, e, NewRuntimeError(e)
	}

	return o, e, nil
}
