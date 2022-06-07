package model

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

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
