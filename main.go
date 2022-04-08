package main

import (
	"certbot-renewer/internal/api"

	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Usage:")
		fmt.Fprintf(w, "# GET - %s/request\r\n", r.Host)
		fmt.Fprintf(w, "Used to request a new SSL certificate for a given domain.\r\n\r\n")

		fmt.Fprintf(w, "# GET - %s/check\r\n", r.Host)
		fmt.Fprintf(w, "Used to fetch the expiration date of a SSL certificate a given domain.\r\n\r\n")

		fmt.Fprintf(w, "# GET - %s/renew\r\n", r.Host)
		fmt.Fprintf(w, "Used to force renew a SSL certificate for a given domain.\r\n")
	})

	http.HandleFunc("/api/request/tls", api.RequestCertificateWithTLS)
	http.HandleFunc("/api/request/http", api.RequestCertificateWithHTTP)
	http.HandleFunc("/api/request/dns", api.RequestCertificateWithDNS)

	if err := http.ListenAndServe(":9022", nil); err != nil {
		log.Printf("listenAndServe failed: %v", err)
	}
}
