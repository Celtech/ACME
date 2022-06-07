package middleware

import (
	"errors"
	"net/http"

	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

var (
	// Authentication
	ErrorUnauthorized    = errors.New("you must be authorized to access this resource")
	ErrorInvalidLogin    = errors.New("the supplied authentication credentials are incorrect")
	ErrorInvalidJWTToken = errors.New("the supplied JWT token was expired or not a valid JWT token")

	ErrorBadPathParameter = errors.New("the supplied path parameter is invalid, did you supply a valid integer?")
)

type ErrorResponse struct {
	Status  int      `json:"status"`
	Errors  []string `json:"errors"`
	Message string   `json:"message"`
}

// ErrorHandler is middleware that enables you to configure error handling from a centralised place via its fluent API.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			switch err.Err {
			case ErrorUnauthorized:
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"errors":  ParseError(err.Err),
					"message": "a valid JWT token is missing from this request",
				})
			case ErrorInvalidLogin:
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"errors":  ParseError(err.Err),
					"message": "Invalid credentials",
				})
			case ErrorInvalidJWTToken:
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"errors":  ParseError(err.Err),
					"message": "Authentication Error",
				})
			case ErrorBadPathParameter:
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusBadRequest,
					"errors":  ParseError(err.Err),
					"message": "Unexpected value for path parameter",
				})
			case gorm.ErrRecordNotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"status":  http.StatusNotFound,
					"errors":  ParseError(errors.New("entity not found")), // we don't like gorms default message
					"message": fmt.Sprintf("The requested entity %s could not be found", c.Param("id")),
				})
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusBadRequest,
					"errors":  ParseError(err.Err),
					"message": "Validation error",
				})
			}
		}
	}
}

func ParseError(errs ...error) []string {
	var out []string
	for _, err := range errs {
		switch typedError := err.(type) {
		case validator.ValidationErrors:
			// if the type is validator.ValidationErrors then it's actually an array of validator.FieldError so we'll
			// loop through each of those and convert them one by one
			for _, e := range typedError {
				out = append(out, parseFieldError(e))
			}
		case *json.UnmarshalTypeError:
			// similarly, if the error is an unmarshalling error we'll parse it into another, more readable string format
			out = append(out, parseMarshallingError(*typedError))
		default:
			out = append(out, err.Error())
		}
	}
	return out
}

func parseFieldError(e validator.FieldError) string {
	// workaround to the fact that the `gt|gtfield=Start` gets passed as an entire tag for some reason
	// https://github.com/go-playground/validator/issues/926
	fieldPrefix := fmt.Sprintf("The field %s", e.Field())
	tag := strings.Split(e.Tag(), "|")[0]
	switch tag {
	case "required_without":
		return fmt.Sprintf("%s is required if %s is not supplied", fieldPrefix, e.Param())
	case "lt", "ltfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be less than %s", fieldPrefix, param)
	case "gt", "gtfield":
		param := e.Param()
		if param == "" {
			param = time.Now().Format(time.RFC3339)
		}
		return fmt.Sprintf("%s must be greater than %s", fieldPrefix, param)
	default:
		// if it's a tag for which we don't have a good format string yet we'll try using the default english translator
		english := en.New()
		translator := ut.New(english, english)
		if translatorInstance, found := translator.GetTranslator("en"); found {
			return e.Translate(translatorInstance)
		} else {
			return fmt.Errorf("%v", e).Error()
		}
	}
}

func parseMarshallingError(e json.UnmarshalTypeError) string {
	return fmt.Sprintf("The field %s must be a %s", e.Field, e.Type.String())
}
