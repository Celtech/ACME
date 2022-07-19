package v1

import (
	"fmt"
	"github.com/Celtech/ACME/internal/queue"
	"github.com/Celtech/ACME/web/middleware"
	"github.com/Celtech/ACME/web/model"
	"net/http"
	"time"

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
// @Success 200 {object} model.APIEnvelopeResponse{data=model.Request}
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
// @Param limit query int false "Number of results per page"
// @Param page query int false "Which page of results to fetch"
// @Param sort query string false "Order of which results appear" Enums(asc, desc)
// @Produce json
// @Success 200 {object} model.APIEnvelopeResponse{data=[]model.Request}
// @Failure 401 {object} middleware.ErrorResponse
// @Router /request [get]
func (requestController RequestController) GetAll(c *gin.Context) {
	pagination := model.GeneratePaginationFromRequest(c)
	requestModel := new(model.Request)
	res, _ := requestModel.GetAll(pagination)

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
// @Success 200 {object} model.APIEnvelopeResponse{data=model.Request}
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

	evt := queue.QueueEvent{
		RequestId:     requestModel.Id,
		Domain:        requestModel.Domain,
		ChallengeType: requestModel.ChallengeType,
		Type:          queue.EVENT_ISSUE,
		Attempt:       1,
		CreatedAt:     time.Now(),
	}

	if err := queue.QueueMgr.Publish(evt); err != nil {
		log.Errorf("error publishing certificate request for domain %s to queue, %v", requestModel.Domain, err)

		c.JSON(400, gin.H{
			"message": fmt.Sprintf("error publishing certificate request for domain %s to queue", requestModel.Domain),
			"error":   err.Error(),
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status": http.StatusCreated,
			"message": fmt.Sprintf(
				"queued certificate request for %s",
				requestModel.Domain,
			),
			"data": requestModel,
		})
	}
}
