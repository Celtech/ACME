package web

import (
	"baker-acme/web/api"

	docs "baker-acme/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func configRoutes(router *gin.Engine) {
	rootGroup := router.Group("/")
	{
		rootRouter(rootGroup)
		apiRouterV1(rootGroup)
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func apiRouterV1(router *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	apiGroup := router.Group("/api/v1")
	{
		apiRequestGroup := apiGroup.Group("/certificate")
		{
			apiRequestGroup.POST("/", api.RequestCertificate)
		}
	}
}
