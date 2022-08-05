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
