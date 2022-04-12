package acme

import (
	"os"
	"time"

	"baker-acme/internal/acme/challenge"
	"baker-acme/internal/acme/challenge/dns01"
	"baker-acme/internal/acme/challenge/http01"
	"baker-acme/internal/acme/challenge/tlsalpn01"

	"github.com/go-acme/lego/v4/providers/dns"
	log "github.com/sirupsen/logrus"

	lego "baker-acme/internal/acme/client"
)

func SetupChallenges(client *lego.Client, challengeType string) error {
	log.Infof("Starting challenge servers")

	switch challengeType {
	case CHALLENGE_TYPE_HTTP:
		if err := client.Challenge.SetHTTP01Provider(setupHTTPProvider()); err != nil {
			return err
		}

	case CHALLENGE_TYPE_TLS:
		if err := client.Challenge.SetTLSALPN01Provider(setupTLSProvider()); err != nil {
			return err
		}

	case CHALLENGE_TYPE_DNS:
		if err := setupDNSProvider(client); err != nil {
			return err
		}
	}

	return nil
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
		log.Infof("HTTP challenge server using proxy header %s", httpProxyHeader)
	}

	log.Infof("HTTP challenge server listening on %s:%s", httpHost, httpPort)

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

	log.Infof("TLS challenge server listening on %s:%s", tlsHost, tlsPort)

	return tlsalpn01.NewProviderServer(tlsHost, tlsPort)
}

func setupDNSProvider(client *lego.Client) error {
	provider, err := dns.NewDNSChallengeProviderByName("dnsmadeeasy")
	if err != nil {
		return err
	}

	return client.Challenge.SetDNS01Provider(provider,
		// dns01.CondOption(len(servers) > 0,
		// 	dns01.AddRecursiveNameservers(dns01.ParseNameservers(ctx.StringSlice("dns.resolvers")))),
		dns01.CondOption(true,
			dns01.DisableCompletePropagationRequirement()),
		dns01.CondOption(true,
			dns01.AddDNSTimeout(time.Duration(60)*time.Second)),
	)
}
