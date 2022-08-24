package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct{}

func (controller IndexController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func (controller IndexController) About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title": "Main website",
	})
}
