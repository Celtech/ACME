package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Celtech/ACME/web/service"
)

// AuthorizeJWT is a middleware for requiring authentication on protected routes
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(ErrorUnauthorized)
			c.Abort()

			return
		}

		tokenString := authHeader[len(BEARER_SCHEMA):]
		valid, err := service.JWTAuthService().ValidateToken(strings.Trim(tokenString, " "))
		if err != nil || !valid {
			c.Error(ErrorInvalidJWTToken)
			c.Abort()

			return
		}
	}
}
