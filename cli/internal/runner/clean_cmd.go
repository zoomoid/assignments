/*
Copyright 2022 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/zoomoid/assignments/v1/internal/config"
)

type cmdCleaner struct {
	*RunnerContext
}

var (
	defaultRecipe = &config.Recipe{
		{
			Command: DefaultBuildProgram,
			Args:    []string{"-C"},
		},
	}
)

// Compile-time type checking of cleaner spec
var _ Cleaner = &cmdCleaner{}

// MakeClean implements the Runner spec in terms of making a singleton exec.Cmd using
// latexmk to cleanup the working directory of the LaTeX compiler
func (c *cmdCleaner) MakeCommand() ([]*exec.Cmd, error) {
	var recipe *config.Recipe
	if c.configuration.Spec.BuildOptions.Cleanup.Command == nil {
		recipe = defaultRecipe
	} else {
		recipe = c.configuration.Spec.BuildOptions.Cleanup.Command.Recipe
		if len(*recipe) == 0 {
			recipe = defaultRecipe
		}
	}

	cmds, err := commandsFromRecipe(recipe, c.TargetDirectory(), c.Filename(), c.Quiet())
	return cmds, err
}

// Run implements the Runner spec in terms of running the cleanup command in shell
func (c *cmdCleaner) Run() error {
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
	log.Debug().Msgf("[runner/clean] Finished cleaning up %s with latexmk", c.TargetDirectory())
	return nil
}
