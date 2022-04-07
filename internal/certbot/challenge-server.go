package certbot

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/providers/dns"

	"github.com/go-acme/lego/v4/lego"
)

func SetupChallenges(client *lego.Client, challengeType string) {
	fmt.Println("Starting HTTP and TLS challenge servers")

	switch challengeType {
	case CHALLENGE_TYPE_HTTP:
		err := client.Challenge.SetHTTP01Provider(setupHTTPProvider())
		if err != nil {
			fmt.Println(err.Error())
		}

	case CHALLENGE_TYPE_TLS:
		errTLS := client.Challenge.SetTLSALPN01Provider(setupTLSProvider())
		if errTLS != nil {
			fmt.Println(errTLS.Error())
		}

	case CHALLENGE_TYPE_DNS:
		setupDNSProvider(client)

	default:
		// TODO: error out and pass up, how did we even get here?
	}
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

func setupDNSProvider(client *lego.Client) {
	provider, err := dns.NewDNSChallengeProviderByName("dnsmadeeasy")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = client.Challenge.SetDNS01Provider(provider,
		// dns01.CondOption(len(servers) > 0,
		// 	dns01.AddRecursiveNameservers(dns01.ParseNameservers(ctx.StringSlice("dns.resolvers")))),
		dns01.CondOption(true,
			dns01.DisableCompletePropagationRequirement()),
		dns01.CondOption(true,
			dns01.AddDNSTimeout(time.Duration(60)*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}
}
