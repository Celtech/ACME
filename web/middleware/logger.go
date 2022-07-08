package middleware

import (
	"bytes"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write I have no idea what this is for if im being honest...
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware is a middleware for logging every http request that comes through the server
func LoggingMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	statusCode := strconv.Itoa(c.Writer.Status())

	log.WithFields(log.Fields{
		"method": c.Request.Method,
		"host":   c.Request.Host,
		"ip":     c.ClientIP(),
		"proto":  c.Request.Proto,
		"status": statusCode,
	}).Infof("Request %s", c.Request.URL)
}
