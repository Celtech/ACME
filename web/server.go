package web

import (
	"baker-acme/web/middleware"
	"fmt"
	"net"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Serve(conf *viper.Viper) *http.Server {
	if os.Getenv("ACME_ENV") == "development" {
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
			conf.GetString("server.host"),
			conf.GetInt("server.port"),
		),
		Handler: router,
	}

	go func() {
		log.Infof(
			"server starting, listening on %s:%d",
			conf.GetString("server.host"),
			conf.GetInt("server.port"),
		)

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
				log.Fatalf("failed to start the server %v", err)
			}
		}
	}()

	return srv
}
