package release

import (
	"os"

	"github.com/rs/zerolog/log"
)

func OpenOrFallbackToStdout(file string) (*os.File, bool) {
	out := os.Stdout
	if file != "" {
		fd, err := os.Open(file)
		if err != nil {
			log.Warn().Err(err).Msg("Falling back to stdout")
			return out, true
		}
		return fd, false
	}
	return out, true
}
