package context

import (
	configFactory "baker-acme/internal/context/config-factory"
	"net/http"

	"go.uber.org/zap"
)

type AppContext struct {
	Logger        *zap.Logger
	HttpWriter    http.ResponseWriter
	ConfigFactory *configFactory.ConfigFactory
}

func NewAppContext(w http.ResponseWriter) *AppContext {
	loggerFactory := newLogger()

	return &AppContext{
		Logger:        loggerFactory,
		HttpWriter:    w,
		ConfigFactory: configFactory.NewConfigFactory(loggerFactory),
	}
}
