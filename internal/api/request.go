package api

import (
	"certbot-renewer/internal/acme"
	"certbot-renewer/internal/domain"
	"fmt"
	"log"
	"net/http"
)

func RequestCertificateWithDNS(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "DNS")

	requestCertificate(w, req, acme.CHALLENGE_TYPE_DNS)
}

func RequestCertificateWithHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "HTTP")

	requestCertificate(w, req, acme.CHALLENGE_TYPE_HTTP)
}

func RequestCertificateWithTLS(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "TLS")

	requestCertificate(w, req, acme.CHALLENGE_TYPE_TLS)
}

func requestCertificate(w http.ResponseWriter, req *http.Request, challengeType string) {
	domainName, err := domain.ParseDomainName(req.Host)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Fprintf(w, "Checking %s", domainName)
	acme.Run(w, domainName, challengeType)
}
