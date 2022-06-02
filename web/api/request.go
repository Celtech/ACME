package api

import (
	"baker-acme/internal/util"
	"fmt"
	"net/http"
	"time"

	"baker-acme/internal/queue"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateCertificateInput struct {
	Domain        string `json:"domain" binding:"required"`
	ChallengeType string `json:"challengeType" binding:"required"`
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func RequestCertificate(c *gin.Context) {
	var input CreateCertificateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainName, err := util.ParseDomainName(input.Domain)
	if err != nil {
		log.Error(err)
	}

	challengeType, err := util.ParseDomainName(input.ChallengeType)
	if err != nil {
		log.Error(err)
	}

	requestCertificate(c, domainName, challengeType)
}

// PongExample godoc
// @Summary pong example
// @Security     ApiKeyAuth
// @Schemes
// @Description do pong
// @Tags idk
// @Accept json
// @Produce json
// @Param	account	body	int  true	"Add account"
// @Success 200 {string} string	"ok"
// @Router /certificate/request [post]
// x-codeSamples file
func requestCertificate(c *gin.Context, domainName string, challengeType string) {
	evt := queue.QueueEvent{
		Domain:        domainName,
		ChallengeType: challengeType,
		Type:          queue.EVENT_REQUEST,
		Attempt:       0,
		CreatedAt:     time.Now(),
	}
	if err := queue.QueueMgr.Publish(evt); err != nil {
		log.Errorf("error publishing certificate request for domain %s to queue, %v", domainName, err)

		c.JSON(400, gin.H{
			"message": fmt.Sprintf("error publishing certificate request for domain %s to queue", domainName),
			"error":   err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf(
				"queued certificate request for %s",
				domainName,
			),
			"request_status": "http://example.com/api/request/1234564634",
		})
	}
}
