package certbot

import (
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
)

func SetupChallenges(client *lego.Client) {
	fmt.Println("Starting HTTP and TLS challenge servers")

	err := client.Challenge.SetHTTP01Provider(setupHTTPProvider())
	if err != nil {
		fmt.Println(err.Error())
	}

	// errTLS := client.Challenge.SetTLSALPN01Provider(setupTLSProvider())
	// if errTLS != nil {
	// 	fmt.Println(errTLS.Error())
	// }
}

func setupHTTPProvider() challenge.Provider {
	httpHost := os.Getenv("HTTP_CHALLENGE_HOST")
	if len(httpHost) == 0 {
		httpHost = "0.0.0.0"
	}

	httpPort := os.Getenv("HTTP_CHALLENGE_PORT")
	if len(httpPort) == 0 {
		httpPort = "5001"
	}

	srv := http01.NewProviderServer(httpHost, httpPort)
	httpProxyHeader := os.Getenv("HTTP_CHALLENGE_PROXY_HEADER")
	if len(httpProxyHeader) > 0 {
		srv.SetProxyHeader(httpProxyHeader)
		fmt.Printf("HTTP challenge server using proxy header %s\n", httpProxyHeader)
	}

	fmt.Printf("HTTP challenge server listening on %s:%s\n", httpHost, httpPort)

	return srv
}

func setupTLSProvider() challenge.Provider {
	tlsHost := os.Getenv("TLS_CHALLENGE_HOST")
	if len(tlsHost) == 0 {
		tlsHost = "0.0.0.0"
	}

	tlsPort := os.Getenv("TLS_CHALLENGE_PORT")
	if len(tlsPort) == 0 {
		tlsPort = "5002"
	}

	fmt.Printf("TLS challenge server listening on %s:%s\n", tlsHost, tlsPort)

	return tlsalpn01.NewProviderServer(tlsHost, tlsPort)
}
