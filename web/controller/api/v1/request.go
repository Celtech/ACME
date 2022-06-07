package v1

import (
	"baker-acme/web/middleware"
	"baker-acme/web/model"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RequestController struct{}

// @BasePath /api/v1

// @Summary Fetch a certificate request
// @Schemes
// @Description Fetch one certificate request by certificate request ID
// @Tags Request
// @Accept json
// @Produce json
// @Param id path int true "Certificate Request ID"
// @Success 200 {object} model.Request
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Router /request/{id} [get]
func (requestController RequestController) GetOne(c *gin.Context) {
	var requestModel = new(model.Request)
	if c.Param("id") == "" {
		c.Error(middleware.ErrorBadPathParameter)
		c.Abort()
		return
	}

	if err := requestModel.GetByID(c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Get one certificate request",
		"data":    requestModel,
	})
}

// @BasePath /api/v1

// @Summary Fetch all certificate requests
// @Schemes
// @Description Fetch all certificate requests
// @Tags Request
// @Accept json
// @Produce json
// @Success 200 {array} model.Request
// @Router /request [get]
func (requestController RequestController) GetAll(c *gin.Context) {
	var requestModel = new(model.Request)
	res, _ := requestModel.GetAll()

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Get all certificate requests",
		"data":    res,
	})
}

// @BasePath /api/v1

// @Summary Create a certificate request
// @Schemes
// @Description Create a certificate request
// @Tags Request
// @Accept json
// @Produce json
// @Param request body model.RequestCreate true "Certificate Request"
// @Success 200 {object} model.Request
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /request [post]
func (requestController RequestController) CreateNew(c *gin.Context) {
	var requestModel = new(model.Request)
	if err := c.ShouldBindJSON(&requestModel); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	requestModel.Status = model.STATUS_PENDING

	if err := requestModel.Save(); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Request POST",
		"data":    requestModel,
	})
}
