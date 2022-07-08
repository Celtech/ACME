package model

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// APIEnvelopeResponse is a generic API Response used only for openapi generation
type APIEnvelopeResponse struct {
	Status  int         `json:"status" example:"201"`
	Message string      `json:"message" example:"use this JWT token as a bearer token to authenticate into the API"`
	Data    interface{} `json:"data"`
}

// Pagination is a struct used to store pagination data for GetAll API requests
type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

// GeneratePaginationFromRequest is a method used to turn extract data from a http request
// into a Pagination struct
func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	limit := 10
	page := 1
	sort := "id asc"
	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "sort":
			sort = queryValue
		}
	}

	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}
