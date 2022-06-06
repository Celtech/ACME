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
	"gorm.io/gorm"
)

func Serve(conf *viper.Viper) *http.Server {
	if os.Getenv("ACME_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware)
	errorHandler(router)
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

func errorHandler(router *gin.Engine) {
	router.Use(middleware.ErrorHandler(
		middleware.Map(gorm.ErrRecordNotFound).
			ToStatusCode(http.StatusNotFound).
			ToResponse(func(c *gin.Context, err error) {
				c.Status(http.StatusNotFound)
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"message": "Entity not found",
					"error":   fmt.Sprintf("The requested entity %s could not be found", c.Param("id")),
				})
			}),
	))

}
