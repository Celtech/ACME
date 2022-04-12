package api

import (
	"baker-acme/internal/acme"
	"baker-acme/internal/util"
	"fmt"
	"net/http"
	"time"

	"baker-acme/internal/queue"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateCertificateInput struct {
	Domain string `json:"domain" binding:"required"`
}

func RequestCertificateWithDNS(c *gin.Context) {
	requestCertificate(c, acme.CHALLENGE_TYPE_DNS)
}

func RequestCertificateWithHTTP(c *gin.Context) {
	requestCertificate(c, acme.CHALLENGE_TYPE_HTTP)
}

func RequestCertificateWithTLS(c *gin.Context) {
	requestCertificate(c, acme.CHALLENGE_TYPE_TLS)
}

func requestCertificate(c *gin.Context, challengeType string) {
	var input CreateCertificateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainName, err := util.ParseDomainName(input.Domain)
	if err != nil {
		log.Error(err)
	}

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
