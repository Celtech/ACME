package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	// Authentication
	ErrorUnauthorized    = errors.New("you must be authorized to access this resource")
	ErrorInvalidLogin    = errors.New("the supplied authentication credentials are incorrect")
	ErrorInvalidJWTToken = errors.New("the supplied JWT token was expired or not a valid JWT token")
)

// ErrorHandler is middleware that enables you to configure error handling from a centralised place via its fluent API.
func ErrorHandler(errMap ...*errorMapping) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		lastErr := context.Errors.Last()
		if lastErr == nil {
			return
		}

		for _, e := range errMap {
			for _, e2 := range e.fromErrors {
				if lastErr.Err == e2 {
					e.toResponse(context, lastErr.Err)
					return
				}
			}
		}
	}
}

type errorMapping struct {
	fromErrors   []error
	toStatusCode int
	toResponse   func(ctx *gin.Context, err error)
}

// ToStatusCode specifies the status code returned to a caller when the error is handled.
func (r *errorMapping) ToStatusCode(statusCode int) *errorMapping {
	r.toStatusCode = statusCode
	r.toResponse = func(ctx *gin.Context, err error) {
		ctx.Status(statusCode)
	}
	return r
}

// ToResponse provides more control over the returned response when an error is matched.
func (r *errorMapping) ToResponse(response func(ctx *gin.Context, err error)) *errorMapping {
	r.toResponse = response
	return r
}

// Map enables you to map errors to a given response status code or response body.
func Map(err ...error) *errorMapping {
	return &errorMapping{
		fromErrors: err,
	}
}
