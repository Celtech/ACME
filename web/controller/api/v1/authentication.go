package v1

import (
	"baker-acme/web/middleware"
	"baker-acme/web/model"
	"baker-acme/web/service"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	// log "github.com/sirupsen/logrus"
)

type AuthenticationController struct{}

func (controller AuthenticationController) Authenticate(c *gin.Context) {
	var userModel = new(model.User)
	if err := c.ShouldBindJSON(&userModel); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	if userModel.Authenticate() {
		token := service.JWTAuthService().GenerateToken(userModel.PublicKey, true)
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Request POST",
			"data":    token,
		})
	} else {
		c.Error(middleware.ErrorInvalidLogin)
		c.Abort()
	}
}
