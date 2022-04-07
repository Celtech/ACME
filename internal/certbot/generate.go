package certbot

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

const filePerm os.FileMode = 0o600
const rootPathWarningMessage = `!!!! HEADS UP !!!!
Your account credentials have been saved in your Let's Encrypt
configuration directory at "%s".
You should make a secure backup of this folder now. This
configuration directory will also contain certificates and
private keys obtained from Let's Encrypt so making regular
backups of this folder is ideal.`

func Run(w http.ResponseWriter, domainName string) {
	accountsStorage := NewAccountsStorage()

	account, client := setup(w, accountsStorage)
	SetupChallenges(client)

	if account.Registration == nil {
		reg, err := register(client)
		if err != nil {
			log.Fatalf("Could not complete registration\n\t%v\n", err)
		}

		account.Registration = reg
		if err = accountsStorage.Save(account); err != nil {
			log.Fatal(err)
		}

		fmt.Printf(rootPathWarningMessage, accountsStorage.GetRootPath())
	}

	certsStorage := NewCertificatesStorage()
	certsStorage.CreateRootFolder()

	cert, err := obtainCertificate(domainName, client)
	if err != nil {
		// Make sure to return a non-zero exit code if ObtainSANCertificate returned at least one error.
		// Due to us not returning partial certificate we can just exit here instead of at the end.
		log.Fatalf("Could not obtain certificates:\n\t%v", err)
	}

	certsStorage.SaveResource(cert)

	// meta := map[string]string{
	// 	renewEnvAccountEmail: account.Email,
	// 	renewEnvCertDomain:   cert.Domain,
	// 	renewEnvCertPath:     certsStorage.GetFileName(cert.Domain, ".crt"),
	// 	renewEnvCertKeyPath:  certsStorage.GetFileName(cert.Domain, ".key"),
	// }

	// return launchHook(ctx.String("run-hook"), meta)
}

func setup(w http.ResponseWriter, accountsStorage *AccountsStorage) (*Account, *lego.Client) {
	privateKey := accountsStorage.GetPrivateKey(acctKeyType)

	var account *Account
	if accountsStorage.ExistsAccountFilePath() {
		account = accountsStorage.LoadAccount(privateKey)
	} else {
		account = &Account{Email: accountsStorage.GetUserID(), key: privateKey}
	}

	client := newClient(w, account, acctKeyType)

	return account, client
}

func register(client *lego.Client) (*registration.Resource, error) {
	return client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func newClient(w http.ResponseWriter, acc registration.User, keyType certcrypto.KeyType) *lego.Client {
	acmeHost := os.Getenv("ACME_HOST")
	if len(acmeHost) == 0 {
		acmeHost = acmeServer
	}

	config := lego.NewConfig(acc)
	config.CADirURL = acmeHost

	config.Certificate = lego.CertificateConfig{
		KeyType: keyType,
		Timeout: time.Duration(30) * time.Second,
	}
	config.UserAgent = "lego-cli/chargeover"
	config.HTTPClient.Timeout = time.Duration(30) * time.Second

	client, err := lego.NewClient(config)
	if err != nil {
		fmt.Printf("Could not create client: %v\n", err)
		os.Exit(1)
	}

	return client
}

func createNonExistingFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0o700)
	} else if err != nil {
		return err
	}
	return nil
}

func obtainCertificate(domainName string, client *lego.Client) (*certificate.Resource, error) {
	fmt.Printf("Requesting certificate for: %s\n", domainName)

	// obtain a certificate, generating a new private key
	request := certificate.ObtainRequest{
		Domains:                        []string{domainName},
		Bundle:                         false,
		MustStaple:                     false,
		AlwaysDeactivateAuthorizations: false,
	}
	return client.Certificate.Obtain(request)

	// read the CSR
	// csr, err := readCSRFile(ctx.String("csr"))
	// if err != nil {
	// 	return nil, err
	// }

	// // obtain a certificate for this CSR
	// return client.Certificate.ObtainForCSR(certificate.ObtainForCSRRequest{
	// 	CSR:                            csr,
	// 	Bundle:                         bundle,
	// 	PreferredChain:                 ctx.String("preferred-chain"),
	// 	AlwaysDeactivateAuthorizations: ctx.Bool("always-deactivate-authorizations"),
	// })
}
