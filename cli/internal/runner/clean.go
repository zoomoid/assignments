package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type cleaner struct {
	*RunnerContext
}

// Compile-time type checking of Runner spec
var _ Runner = &cleaner{}

// MakeClean implements the Runner spec in terms of making a singleton exec.Cmd using
// latexmk to cleanup the working directory of the LaTeX compiler
func (c *cleaner) MakeCommand() ([]*exec.Cmd, error) {
	out := &bytes.Buffer{}

	cmd := exec.Command(defaultProgram, "-C")

	if c.quiet {
		cmd.Stdout = out
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Dir = c.TargetDirectory()
	return []*exec.Cmd{cmd}, nil
}

// Run implements the Runner spec in terms of running the cleanup command in shell
func (c *cleaner) Run() error {
	log.Debug().Msgf("[runner/clean] Cleaning up %s using latexmk", c.TargetDirectory())

	cmds, err := c.MakeCommand()
	if err != nil {
		return err
	}

	c.Commands = cmds

	for i, cmd := range c.Commands {
		if cmd == nil {
			return fmt.Errorf("command %d is nil", i)
		}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	log.Debug().Msgf("[runner/clean] Finished cleanup %s with latexmk", c.TargetDirectory())
	return nil
}
