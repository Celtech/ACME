package web

import (
	"baker-acme/internal/context"
	"baker-acme/web/middleware"
	"fmt"
	"net"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func StartServer(appContext *context.AppContext) *http.Server {
	if appContext.ConfigFactory.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware)

	configRoutes(router)

	srv := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%d",
			appContext.ConfigFactory.Server.Host,
			appContext.ConfigFactory.Server.Port,
		),
		Handler: router,
	}

	go func() {
		log.Infof(
			"server starting, listening on %s:%d",
			appContext.ConfigFactory.Server.Host,
			appContext.ConfigFactory.Server.Port,
		)

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
				log.Fatalf("failed to start the server %v", err)
			}
		}
	}()

	return srv
}
