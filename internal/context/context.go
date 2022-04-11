package context

import (
	configFactory "baker-acme/internal/context/config-factory"
	"net/http"
)

type AppContext struct {
	HttpWriter    http.ResponseWriter
	ConfigFactory *configFactory.ConfigFactory
}

func NewAppContext(w http.ResponseWriter) *AppContext {
	return &AppContext{
		HttpWriter:    w,
		ConfigFactory: configFactory.NewConfigFactory(),
	}
}
