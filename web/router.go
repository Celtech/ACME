package web

import (
	"baker-acme/web/api"

	"github.com/gin-gonic/gin"
)

func configRoutes(router *gin.Engine) {
	rootGroup := router.Group("/")
	{
		rootRouter(rootGroup)
		apiRouter(rootGroup)
	}

}

func rootRouter(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "welcome",
		})
	})
}

func apiRouter(router *gin.RouterGroup) {
	apiGroup := router.Group("/api")
	{
		apiRequestGroup := apiGroup.Group("/request")
		{
			apiRequestGroup.POST("/tls", api.RequestCertificateWithTLS)
			apiRequestGroup.POST("/http", api.RequestCertificateWithHTTP)
			apiRequestGroup.POST("/dns", api.RequestCertificateWithDNS)
		}
	}
}
