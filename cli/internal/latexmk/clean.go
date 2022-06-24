package latexmk

import (
	"bytes"
	"os"
	"os/exec"
)

type cleaner struct {
	*RunnerContext
}

func (c *cleaner) makeCleanCmd() error {
	out := &bytes.Buffer{}

	cmd := exec.Command(defaultProgram, "-C")

	if c.Quiet {
		cmd.Stdout = out
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Dir = c.TargetDirectory
	c.Commands = []*exec.Cmd{cmd}
	return nil
}

var _ Runner = &cleaner{}

func (c *cleaner) Run() error {
	c.Logger.Debug("[runner/clean] Cleaning up using latexmk", "pwd", c.Root)

	err := c.makeCleanCmd()
	if err != nil {
		return err
	}
	for _, cmd := range c.Commands {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	c.Logger.Debug("[runner/clean] Finished cleanup with latexmk", "directory", c.TargetDirectory)
	return nil
}
