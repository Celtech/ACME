package main

import (
	"certbot-renewer/internal/certbot"
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

func checkCertificate(w http.ResponseWriter, req *http.Request) {
	domainName := domain.ParseDomainName(req.Host)
	fmt.Fprintf(w, fmt.Sprintf("Checking %s", domainName))

	certbot.Run(w, domainName)
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
	http.HandleFunc("/request", checkCertificate)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `                                                                                  
 ██████ ██   ██  █████  ██████   ██████  ███████  ██████  ██    ██ ███████ ██████  
██      ██   ██ ██   ██ ██   ██ ██       ██      ██    ██ ██    ██ ██      ██   ██ 
██      ███████ ███████ ██████  ██   ███ █████   ██    ██ ██    ██ █████   ██████  
██      ██   ██ ██   ██ ██   ██ ██    ██ ██      ██    ██  ██  ██  ██      ██   ██ 
 ██████ ██   ██ ██   ██ ██   ██  ██████  ███████  ██████    ████   ███████ ██   ██`)
		fmt.Fprintln(w, "┌────────────┬───────────────────┬───────────────────────────────────────────────┐")
		fmt.Fprintln(w, fmt.Sprintf("│ Status: OK │ Uptime: %s │                                               │", uptime()))
		fmt.Fprintln(w, "└────────────┴───────────────────┴───────────────────────────────────────────────┘\r\n")

		fmt.Fprintln(w, "Usage:")
		fmt.Fprintln(w, fmt.Sprintf("# GET - %s://%s/request", req(*r).Protocol(), r.Host))
		fmt.Fprintln(w, "Used to request a new SSL certificate for a given domain.\r\n")

		fmt.Fprintln(w, fmt.Sprintf("# GET - %s://%s/check", req(*r).Protocol(), r.Host))
		fmt.Fprintln(w, "Used to fetch the expiration date of a SSL certificate a given domain.\r\n")

		fmt.Fprintln(w, fmt.Sprintf("# GET - %s://%s/renew", req(*r).Protocol(), r.Host))
		fmt.Fprintln(w, "Used to force renew a SSL certificate for a given domain.\r\n")

	})

	if err := http.ListenAndServe(":9022", nil); err != nil {
		log.Printf("listenAndServe failed: %v", err)
	}
}
