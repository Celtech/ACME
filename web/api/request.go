package api

import (
	"baker-acme/internal/acme"
	"baker-acme/internal/util"
	"fmt"
	"net/http"

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

	if err := acme.Run(domainName, challengeType); err != nil {
		log.Errorf("Error issuing certificate for %s\r\n%v", domainName, err)

		c.JSON(400, gin.H{
			"message": fmt.Sprintf("Error issuing certificate for %s", domainName),
			"error":   err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Certificate issued for %s", domainName),
		})
	}
}
