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
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/zoomoid/assignments/v1/internal/util"
)

type globCleaner struct {
	*RunnerContext
}

var (
	DefaultPatterns = []string{
		"*.aux",
		"*.bbl",
		"*.blg",
		"*.idx",
		"*.ind",
		"*.lof",
		"*.lot",
		"*.out",
		"*.toc",
		"*.acn",
		"*.acr",
		"*.alg",
		"*.glg",
		"*.glo",
		"*.gls",
		"*.fls",
		"*.log",
		"*.fdb_latexmk",
		"*.snm",
		"*.synctex(busy)",
		"*.synctex.gz(busy)",
		"*.nav",
		"*.vrb",
	}
)

var _ Cleaner = &globCleaner{}

func (c *globCleaner) MakeCommand() ([]*exec.Cmd, error) {
	// globCleaner does not actually run anything in a shell, so return an
	// empty slice of commands and implement cleaning in Run()
	return []*exec.Cmd{}, nil
}

func (c *globCleaner) Run() error {
	log.Debug().Msg("[runner/clean] Cleaning up using glob patterns")
	globOptions := c.configuration.Spec.BuildOptions.Cleanup.Glob

	patterns := globOptions.Patterns
	recursive := globOptions.Recursive
	if len(patterns) == 0 {
		patterns = DefaultPatterns
	}

	visitors, errs := util.ExpandPaths(c.TargetDirectory(), patterns, recursive)
	elist := util.NewErrorList(errs)
	if elist != nil {
		return elist
	}

	for _, v := range visitors {
		v.Visit(func(path string) error {
			return os.Remove(path)
		})
	}
	log.Debug().Msgf("[runner/clean] Finished cleaning up %s with glob patterns", c.TargetDirectory())
	return nil
}
