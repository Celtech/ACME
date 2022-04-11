package api

import (
	"baker-acme/internal/acme"
	"baker-acme/internal/util"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func RequestCertificateWithDNS(w http.ResponseWriter, req *http.Request) {

	requestCertificate(w, req, acme.CHALLENGE_TYPE_DNS)
}

func RequestCertificateWithHTTP(w http.ResponseWriter, req *http.Request) {
	requestCertificate(w, req, acme.CHALLENGE_TYPE_HTTP)
}

func RequestCertificateWithTLS(w http.ResponseWriter, req *http.Request) {
	requestCertificate(w, req, acme.CHALLENGE_TYPE_TLS)
}

func requestCertificate(w http.ResponseWriter, req *http.Request, challengeType string) {
	domainName, err := util.ParseDomainName(req.Host)
	if err != nil {
		log.Error(err)
	}

	if err := acme.Run(domainName, challengeType); err != nil {
		log.Errorf("Error issuing certificate for %s\r\n%v", domainName, err)
		fmt.Fprintf(w, "Error issuing certificate for %s\r\n%v", domainName, err)
	} else {
		fmt.Fprintf(w, "Certificate issued for %s", domainName)
	}
}
