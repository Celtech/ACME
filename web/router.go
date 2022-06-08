package web

import (
	"baker-acme/web/controller"
	v1 "baker-acme/web/controller/api/v1"
	"baker-acme/web/middleware"

	"github.com/gin-gonic/gin"
)

func configRoutes(router *gin.Engine) {
	rootRouter(router)
	apiV1Router(router)
}

func rootRouter(router *gin.Engine) {
	rootController := new(controller.RootController)

	router.StaticFile("/openapi", "./docs/swagger.json")
	router.GET("/ping", rootController.Ping)
}

func apiV1Router(router *gin.Engine) {
	v1Group := router.Group("/api/v1")
	{
		authGroup := v1Group.Group("token")
		{
			authController := new(v1.AuthenticationController)
			authGroup.POST("/", authController.Authenticate)
		}

		requestGroup := v1Group.Group("request")
		{
			requestController := new(v1.RequestController)
			requestGroup.Use(middleware.AuthorizeJWT())
			requestGroup.GET("/", requestController.GetAll)
			requestGroup.GET("/:id", requestController.GetOne)
			requestGroup.POST("/", requestController.CreateNew)
		}
	}
}
