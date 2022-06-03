package web

import (
	"baker-acme/web/api"

	"github.com/gin-gonic/gin"
)

func configRoutes(router *gin.Engine) {
	rootGroup := router.Group("/")
	{
		rootRouter(rootGroup)
		apiRouterV1(rootGroup)
	}

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.LoadHTMLFiles("./docs/*")
	router.StaticFile("/swagger", "./docs/swagger.json")
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

func apiRouterV1(router *gin.RouterGroup) {
	apiGroup := router.Group("/api/v1")
	{
		apiRequestGroup := apiGroup.Group("/certificate")
		{
			apiRequestGroup.POST("/", api.RequestCertificate)
		}
	}
}
