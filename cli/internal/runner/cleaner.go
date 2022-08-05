package runner

import (
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/zoomoid/assignments/v1/internal/config"
)

type Cleaner interface {
	Runner
}

var (
	DefaultCleaner *config.CleanupOptions = &config.CleanupOptions{
		Glob: &config.CleanupGlobOptions{
			Recursive: false,
			Patterns:  DefaultPatterns,
		},
	}
)

// dummyCleaner is a cleaner implementation that does effectively nothing but implements
// the required interface
type dummyCleaner struct{}

var _ Cleaner = &dummyCleaner{}

func (c *dummyCleaner) MakeCommand() ([]*exec.Cmd, error) {
	return nil, nil
}

func (c *dummyCleaner) Run() error {
	log.Debug().Msg("Skipping cleanup because spec.buildOptions.cleanup is nil")
	return nil
}
