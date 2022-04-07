package main

import (
	"certbot-renewer/internal/api"
	"certbot-renewer/internal/domain"

	"time"

	"fmt"
	"log"
	"net/http"
)

var startTime time.Time

type Timespan time.Duration

func (t Timespan) Format(format string) string {
	z := time.Unix(0, 0).UTC()
	return z.Add(time.Duration(t)).Format(format)
}

func uptime() string {
	return Timespan(time.Since(startTime).Round(time.Second)).Format("15h04m05s")
}

func init() {
	startTime = time.Now()
}

type req http.Request

func (t req) Protocol() string {
	scheme := "http"
	if t.TLS != nil {
		scheme = "https"
	}

	return scheme
}

func main() {
	log.Println("Starting server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		domainName, err := domain.ParseDomainName(r.Host)
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Fprintln(w, `                                                                                  
  ██████ ██   ██  █████  ██████   ██████  ███████  ██████  ██    ██ ███████ ██████  
 ██      ██   ██ ██   ██ ██   ██ ██       ██      ██    ██ ██    ██ ██      ██   ██ 
 ██      ███████ ███████ ██████  ██   ███ █████   ██    ██ ██    ██ █████   ██████  
 ██      ██   ██ ██   ██ ██   ██ ██    ██ ██      ██    ██  ██  ██  ██      ██   ██ 
  ██████ ██   ██ ██   ██ ██   ██  ██████  ███████  ██████    ████   ███████ ██   ██`)
		fmt.Fprintf(w, "┌─                                                                                ─┐\r\n")
		fmt.Fprintf(w, "  Status: OK │ Uptime: %s │ Host: %s \r\n", uptime(), domainName)
		fmt.Fprintf(w, "└─                                                                                ─┘\r\n\r\n")

		fmt.Fprintln(w, "Usage:")
		fmt.Fprintf(w, "# GET - %s://%s/request\r\n", req(*r).Protocol(), r.Host)
		fmt.Fprintf(w, "Used to request a new SSL certificate for a given domain.\r\n\r\n")

		fmt.Fprintf(w, "# GET - %s://%s/check\r\n", req(*r).Protocol(), r.Host)
		fmt.Fprintf(w, "Used to fetch the expiration date of a SSL certificate a given domain.\r\n\r\n")

		fmt.Fprintf(w, "# GET - %s://%s/renew\r\n", req(*r).Protocol(), r.Host)
		fmt.Fprintf(w, "Used to force renew a SSL certificate for a given domain.\r\n")
	})

	http.HandleFunc("/api/request/tls", api.RequestCertificateWithTLS)
	http.HandleFunc("/api/request/http", api.RequestCertificateWithHTTP)
	http.HandleFunc("/api/request/dns", api.RequestCertificateWithDNS)

	if err := http.ListenAndServe(":9022", nil); err != nil {
		log.Printf("listenAndServe failed: %v", err)
	}
}
