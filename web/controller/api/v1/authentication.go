package v1

import (
	"baker-acme/web/middleware"
	"baker-acme/web/model"
	"baker-acme/web/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct{}

func (controller AuthenticationController) Authenticate(c *gin.Context) {
	var userModel = new(model.User)
	if err := c.ShouldBindJSON(&userModel); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if !userModel.Authenticate() {
		c.Error(middleware.ErrorInvalidLogin)
		c.Abort()
		return
	}

	token := service.JWTAuthService().GenerateToken(userModel.PublicKey, true)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "use this JWT token as a bearer token to authenticate into the API",
		"data":    token,
	})
}
