package v1

import (
	"github.com/Celtech/ACME/web/middleware"
	"github.com/Celtech/ACME/web/model"
	"github.com/Celtech/ACME/web/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct{}

// @BasePath /api/v1

// @Summary Create JWT token
// @Schemes
// @Description Create JWT token from a users email and password.
// @Description This token is used to authenticate to the rest of the API
// @Tags Token
// @Accept json
// @Produce json
// @Param request body model.User true "Token Request"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /token [post]
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

	token := service.JWTAuthService().GenerateToken(userModel.Email, true)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "use this JWT token as a bearer token to authenticate into the API",
		"data":    token,
	})
}
