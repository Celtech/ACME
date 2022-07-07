package acme

import (
	"baker-acme/config"
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
	httpHost := config.GetConfig().GetString("acme.http.host")
	if len(httpHost) == 0 {
		httpHost = "0.0.0.0"
	}

	httpPort := config.GetConfig().GetString("acme.http.port")
	if len(httpPort) == 0 {
		httpPort = "80"
	}

	srv := http01.NewProviderServer(httpHost, httpPort)
	httpProxyHeader := config.GetConfig().GetString("acme.http.proxyHeader")
	if len(httpProxyHeader) > 0 {
		srv.SetProxyHeader(httpProxyHeader)
		log.Infof("HTTP challenge server using proxy header %s", httpProxyHeader)
	}

	log.Infof("HTTP challenge server listening on %s:%s", httpHost, httpPort)

	return srv
}

func setupTLSProvider() challenge.Provider {
	tlsHost := config.GetConfig().GetString("acme.tls.host")
	if len(tlsHost) == 0 {
		tlsHost = "0.0.0.0"
	}

	tlsPort := config.GetConfig().GetString("acme.tls.port")
	if len(tlsPort) == 0 {
		tlsPort = "443"
	}

	srv := tlsalpn01.NewProviderServer(tlsHost, tlsPort)

	log.Infof("TLS challenge server listening on %s:%s", tlsHost, tlsPort)

	return srv
}

func setupDNSProvider(client *lego.Client) error {
	dnsProvider := config.GetConfig().GetString("acme.dns.provider")
	if len(dnsProvider) == 0 {
		dnsProvider = "dnsmadeeasy"
	}

	dnsTimeout := config.GetConfig().GetInt("acme.dns.timeout")
	if dnsTimeout == 0 {
		dnsTimeout = 60
	}

	provider, err := dns.NewDNSChallengeProviderByName(dnsProvider)
	if err != nil {
		return err
	}

	return client.Challenge.SetDNS01Provider(provider,
		dns01.CondOption(true,
			dns01.DisableCompletePropagationRequirement()),
		dns01.CondOption(true,
			dns01.AddDNSTimeout(time.Duration(dnsTimeout)*time.Second)),
	)
}
