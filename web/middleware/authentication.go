package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"baker-acme/web/service"
)

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
		token, err := service.JWTAuthService().ValidateToken(strings.Trim(tokenString, " "))
		if err != nil || !token.Valid {
			c.Error(ErrorInvalidJWTToken)
			c.Abort()

			return
		}
	}
}
