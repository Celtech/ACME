package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "baker-acme/web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	ErrorNotFound = fmt.Errorf("this is an error")
)

type ValidationError struct{}

func (e *ValidationError) Error() string {
	return "Invalid request"
}

func TestMapSimpleErrorToStatusCode(t *testing.T) {
	// Arrange
	router := gin.Default()
	router.Use(
		ErrorHandler(
			Map(ErrorNotFound).ToStatusCode(http.StatusNotFound),
		))

	// Act
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(ErrorNotFound)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusNotFound)
}

func TestMapErrorStructToStatusCode(t *testing.T) {
	// Arrange
	router := gin.Default()
	router.Use(
		ErrorHandler(
			Map(&ValidationError{}).ToStatusCode(http.StatusBadRequest),
		))

	// Act
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(&ValidationError{})
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusBadRequest)
}

func TestMapErrorResponseFunc(t *testing.T) {
	// Arrange
	router := gin.Default()
	router.Use(
		ErrorHandler(
			Map(ErrorNotFound).ToResponse(func(c *gin.Context, err error) {
				c.Status(http.StatusNotFound)
				c.Writer.Write([]byte(err.Error()))
			}),
		))

	// Act
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(ErrorNotFound)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	// Assert
	assert.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
	assert.Equal(t, ErrorNotFound.Error(), recorder.Body.String())
}
