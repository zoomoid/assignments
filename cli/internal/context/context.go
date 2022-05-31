package context

import (
	config "github.com/zoomoid/assignments/v1/internal/config"
	zap "go.uber.org/zap"
)

type AppContext struct {
	Logger        *zap.SugaredLogger
	Cwd           string
	Configuration *config.Configuration
}
