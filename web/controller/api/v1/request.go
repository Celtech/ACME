package v1

import (
	"baker-acme/web/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestController struct{}

var requestModel = new(model.Request)

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
	if c.Param("id") != "" {
		cert, err := requestModel.GetByID(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve user", "error": err})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User founded!", "user": cert})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	c.Abort()
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
	c.JSON(http.StatusOK, gin.H{"message": "Request GETALL"})
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
	c.JSON(http.StatusOK, gin.H{"message": "Request POST"})
}
