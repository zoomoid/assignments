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

package ci

import (
	"os"

	"github.com/rs/zerolog/log"
)

func OpenOrFallbackToStdout(file string) (*os.File, bool) {
	out := os.Stdout
	if file != "" {
		fd, err := os.Create(file)
		if err != nil {
			log.Warn().Err(err).Msg("Falling back to stdout")
			return out, true
		}
		return fd, false
	}
	return out, true
}
