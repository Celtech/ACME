package web

import (
	"baker-acme/internal/context"
	"baker-acme/web/api"
	"fmt"

	"net"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

func StartServer(appContext *context.AppContext, wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{
		Addr: fmt.Sprintf(
			"%s:%d",
			appContext.ConfigFactory.Server.Host,
			appContext.ConfigFactory.Server.Port,
		),
	}

	router()

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		appContext.Logger.Info(
			fmt.Sprintf(
				"server starting, listening on %s:%d",
				appContext.ConfigFactory.Server.Host,
				appContext.ConfigFactory.Server.Port,
			),
		)

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
				appContext.Logger.Fatal("failed to start the server",
					zap.Error(err),
				)
			}
		}
	}()

	return srv
}

func router() {
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
}
