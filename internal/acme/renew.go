package acme

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/Celtech/ACME/config"
	"github.com/Celtech/ACME/internal/acme/certificate"
	lego "github.com/Celtech/ACME/internal/acme/client"
	certcrypto "github.com/Celtech/ACME/internal/acme/crypto"
	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

func Renew(domains []string, challengeType string) error {
	accountsStorage, err := NewAccountsStorage()
	if err != nil {
		return err
	}

	account, client, err := setup(accountsStorage)
	if err != nil {
		return err
	}

	if err := SetupChallenges(client, challengeType); err != nil {
		return err
	}

	if account.Registration == nil {
		return fmt.Errorf("Account %s is not registered. Use 'run' to register a new account.\n", account.Email)
	}

	certsStorage := NewCertificatesStorage()
	//// Domains
	return renewForDomains(domains, client, certsStorage, false)
}

func renewForDomains(domains []string, client *lego.Client, certsStorage *CertificatesStorage, bundle bool) error {
	domain := domains[0]

	// load the cert resource from files.
	// We store the certificate, private key and metadata in different files
	// as web servers would not be able to work with a combined file.
	certificates, err := certsStorage.ReadCertificate(domain, ".crt")
	if err != nil {
		return fmt.Errorf("Error while loading the certificate for domain %s\n\t%v", domain, err)
	}

	cert := certificates[0]

	renewalDays := config.GetConfig().GetInt("acme.renewal.days")
	if renewalDays == 0 {
		renewalDays = 30
	}

	if !needRenewal(cert, domain, renewalDays) {
		return nil
	}

	// This is just meant to be informal for the user.
	timeLeft := cert.NotAfter.Sub(time.Now().UTC())
	log.Infof("[%s] acme: Trying renewal with %d hours remaining", domain, int(timeLeft.Hours()))

	certDomains := certcrypto.ExtractDomains(cert)

	var privateKey crypto.PrivateKey

	reusePrivateKey := config.GetConfig().GetBool("acme.renewal.reusePrivateKey")
	if reusePrivateKey {
		keyBytes, errR := certsStorage.ReadFile(domain, ".key")
		if errR != nil {
			return fmt.Errorf("Error while loading the private key for domain %s\n\t%v", domain, errR)
		}

		privateKey, errR = certcrypto.ParsePEMPrivateKey(keyBytes)
		if errR != nil {
			return errR
		}
	}

	// https://github.com/certbot/certbot/blob/284023a1b7672be2bd4018dd7623b3b92197d4b0/certbot/certbot/_internal/renewal.py#L435-L440
	if !isatty.IsTerminal(os.Stdout.Fd()) && !config.GetConfig().GetBool("acme.renewal.noRandomDelay") {
		// https://github.com/certbot/certbot/blob/284023a1b7672be2bd4018dd7623b3b92197d4b0/certbot/certbot/_internal/renewal.py#L472
		const jitter = 8 * time.Minute
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		sleepTime := time.Duration(rnd.Int63n(int64(jitter)))

		log.Infof("renewal: random delay of %s", sleepTime)
		time.Sleep(sleepTime)
	}

	request := certificate.ObtainRequest{
		Domains:                        merge(certDomains, domains),
		Bundle:                         bundle,
		PrivateKey:                     privateKey,
		MustStaple:                     config.GetConfig().GetBool("acme.renewal.mustStaple"),
		PreferredChain:                 "",
		AlwaysDeactivateAuthorizations: config.GetConfig().GetBool("acme.renewal.alwaysDeactivateAuthorizations"),
	}
	certRes, err := client.Certificate.Obtain(request)
	if err != nil {
		return fmt.Errorf("could not obtain certificates:\n\t%v", err)
	}

	// TODO: call a renewal hook here to trigger some action to happen
	return certsStorage.SaveResource(certRes)
}

func needRenewal(x509Cert *x509.Certificate, domain string, days int) bool {
	if x509Cert.IsCA {
		log.Fatalf("[%s] Certificate bundle starts with a CA certificate", domain)
	}

	if days >= 0 {
		notAfter := int(time.Until(x509Cert.NotAfter).Hours() / 24.0)
		if notAfter > days {
			log.Infof("[%s] The certificate expires in %d days, the number of days defined to perform the renewal is %d: no renewal.",
				domain, notAfter, days)
			return false
		}
	}

	return true
}

func merge(prevDomains, nextDomains []string) []string {
	for _, next := range nextDomains {
		var found bool
		for _, prev := range prevDomains {
			if prev == next {
				found = true
				break
			}
		}
		if !found {
			prevDomains = append(prevDomains, next)
		}
	}
	return prevDomains
}
