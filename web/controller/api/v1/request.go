package v1

import (
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
// @Success 200 {string} Helloworld
// @Success 400 {string} MissingID
// @Success 404 {string} NotFound
// @Router /request/{id} [get]
func (requestController RequestController) GetOne(c *gin.Context) {
	var requestModel = new(model.Request)
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Something went wrong",
			"error":   "The provided ID is invalid",
		})
		c.Abort()
		return
	}

	if err := requestModel.GetByID(c.Param("id")); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Found",
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
// @Success 200 {string} Helloworld
// @Router /request [get]
func (requestController RequestController) GetAll(c *gin.Context) {
	var requestModel = new(model.Request)
	res, err := requestModel.GetAll()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Request GET ALL",
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
// @Param request body model.Request true "Certificate Request"
// @Success 200 {string} Helloworld
// @Success 400 {string} MissingBody
// @Success 422 {string} ValidationError
// @Router /request [post]
func (requestController RequestController) CreateNew(c *gin.Context) {
	var requestModel = new(model.Request)
	if err := c.ShouldBindJSON(&requestModel); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
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
