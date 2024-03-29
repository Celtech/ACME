package acme

import (
	"fmt"
	acmeConfig "github.com/Celtech/ACME/config"
	"os"
	"strings"
	"time"

	"github.com/Celtech/ACME/internal/acme/certificate"
	lego "github.com/Celtech/ACME/internal/acme/client"
	"github.com/Celtech/ACME/internal/acme/crypto"
	"github.com/Celtech/ACME/internal/acme/registration"

	log "github.com/sirupsen/logrus"
)

const (
	CHALLENGE_TYPE_TLS  = "challenge-tls"
	CHALLENGE_TYPE_HTTP = "challenge-http"
	CHALLENGE_TYPE_DNS  = "challenge-dns"
)

const filePerm os.FileMode = 0o666
const rootPathWarningMessage = `!!!! HEADS UP !!!!
Your account credentials have been saved in your Let's Encrypt
configuration directory at "%s".
You should make a secure backup of this folder now. This
configuration directory will also contain certificates and
private keys obtained from Let's Encrypt so making regular
backups of this folder is ideal.`

func Run(domainName string, challengeType string) error {
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
		reg, err := register(client)
		if err != nil {
			return fmt.Errorf(
				"could not complete registration\n\t%v", err,
			)
		}

		account.Registration = reg
		if err = accountsStorage.Save(account); err != nil {
			return err
		}

		log.Infof(rootPathWarningMessage, accountsStorage.GetRootPath())
	}

	certsStorage := NewCertificatesStorage()
	if err := certsStorage.CreateRootFolder(); err != nil {
		return err
	}

	cert, err := obtainCertificate(domainName, client)
	if err != nil {
		return fmt.Errorf(
			"could not obtain certificates:\n\t%v", err,
		)
	}

	return certsStorage.SaveResource(cert)
}

func setup(accountsStorage *AccountsStorage) (*Account, *lego.Client, error) {
	privateKey, err := accountsStorage.GetPrivateKey(acctKeyType)
	if err != nil {
		return nil, nil, err
	}

	var account *Account
	res, err := accountsStorage.ExistsAccountFilePath()
	if err != nil {
		return nil, nil, err
	} else if res {
		account = accountsStorage.LoadAccount(privateKey)
	} else {
		account = &Account{Email: accountsStorage.GetUserID(), key: privateKey}
	}

	client, err := newClient(account, acctKeyType)
	if err != nil {
		return nil, nil, err
	}

	return account, client, nil
}

func register(client *lego.Client) (*registration.Resource, error) {
	return client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func newClient(acc registration.User, keyType crypto.KeyType) (*lego.Client, error) {
	acmeHost := acmeConfig.GetConfig().GetString("acme.host")
	if len(acmeHost) == 0 {
		acmeHost = acmeServer
	}

	acmeUserAgent := acmeConfig.GetConfig().GetString("acme.userAgent")
	if len(acmeUserAgent) == 0 {
		acmeUserAgent = "lego-cli/chargeover"
	}

	clientTimeout := acmeConfig.GetConfig().GetInt("acme.clientTimeout")
	if clientTimeout == 0 {
		clientTimeout = 60
	}

	config := lego.NewConfig(acc)
	config.CADirURL = acmeHost

	config.Certificate = lego.CertificateConfig{
		KeyType: keyType,
		Timeout: time.Duration(clientTimeout) * time.Second,
	}
	config.UserAgent = acmeUserAgent
	config.HTTPClient.Timeout = time.Duration(clientTimeout) * time.Second

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf(
			"could not create client: %v", err,
		)
	}

	return client, nil
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
	log.Infof("requesting certificate for: %s", domainName)

	// Support for wildcard certificates
	var domainNames []string
	if strings.HasPrefix(domainName, "*.") {
		domainNames = append(domainNames, domainName)
		domainNames = append(domainNames, strings.TrimPrefix(domainName, "*."))
	} else {
		domainNames = append(domainNames, domainName)
	}

	// obtain a certificate, generating a new private key
	request := certificate.ObtainRequest{
		Domains:                        domainNames,
		Bundle:                         false,
		MustStaple:                     false,
		AlwaysDeactivateAuthorizations: false,
	}
	return client.Certificate.Obtain(request)
}
