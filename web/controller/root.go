package controller

import (
	"github.com/gin-gonic/gin"
)

type RootController struct{}

// Check that the server is alive and healthy
func (rootController RootController) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
