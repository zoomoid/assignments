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
