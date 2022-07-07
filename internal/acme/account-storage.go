package acme

import (
	acmeConfig "baker-acme/config"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	lego "baker-acme/internal/acme/client"
	certcrypto "baker-acme/internal/acme/crypto"
	"baker-acme/internal/acme/registration"
	log "github.com/sirupsen/logrus"
)

const (
	baseAccountsRootFolderName = "accounts"
	baseKeysFolderName         = "keys"
	accountFileName            = "account.json"
	acmeServer                 = "https://acme-v02.api.letsencrypt.org/directory"
)

// AccountsStorage A storage for account data.
//
// rootPath:
//
//     ./.lego/accounts/
//          │      └── root accounts directory
//          └── "path" option
//
// rootUserPath:
//
//     ./.lego/accounts/localhost_14000/example@chargeover.com/
//          │      │             │             └── userID ("email" option)
//          │      │             └── CA server ("server" option)
//          │      └── root accounts directory
//          └── "path" option
//
// keysPath:
//
//     ./.lego/accounts/localhost_14000/example@chargeover.com/keys/
//          │      │             │             │           └── root keys directory
//          │      │             │             └── userID ("email" option)
//          │      │             └── CA server ("server" option)
//          │      └── root accounts directory
//          └── "path" option
//
// accountFilePath:
//
//     ./.lego/accounts/localhost_14000/example@chargeover.com/account.json
//          │      │             │             │             └── account file
//          │      │             │             └── userID ("email" option)
//          │      │             └── CA server ("server" option)
//          │      └── root accounts directory
//          └── "path" option
//
type AccountsStorage struct {
	userID          string
	rootPath        string
	rootUserPath    string
	keysPath        string
	accountFilePath string
}

// NewAccountsStorage Creates a new AccountsStorage.
func NewAccountsStorage() (*AccountsStorage, error) {
	// TODO: move to account struct? Currently MUST pass email.
	email := acmeConfig.GetConfig().GetString("acme.email")
	if len(email) == 0 {
		return nil, errors.New("you must set the `acme.email` config key in your configuration file")
	}

	acmeHost := acmeConfig.GetConfig().GetString("acme.host")
	if len(acmeHost) == 0 {
		acmeHost = acmeServer
	}
	serverURL, err := url.Parse(acmeHost)
	if err != nil {
		return nil, err
	}

	basePath := acmeConfig.GetConfig().GetString("acme.dataPath")
	if len(basePath) == 0 {
		basePath = "/data"
	}

	rootPath := filepath.Join(basePath, baseAccountsRootFolderName)
	serverPath := strings.NewReplacer(":", "_", "/", string(os.PathSeparator)).Replace(serverURL.Host)
	accountsPath := filepath.Join(rootPath, serverPath)
	rootUserPath := filepath.Join(accountsPath, email)

	return &AccountsStorage{
		userID:          email,
		rootPath:        rootPath,
		rootUserPath:    rootUserPath,
		keysPath:        filepath.Join(rootUserPath, baseKeysFolderName),
		accountFilePath: filepath.Join(rootUserPath, accountFileName),
	}, nil
}

func (s *AccountsStorage) ExistsAccountFilePath() (bool, error) {
	accountFile := filepath.Join(s.rootUserPath, accountFileName)
	if _, err := os.Stat(accountFile); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *AccountsStorage) GetRootPath() string {
	return s.rootPath
}

func (s *AccountsStorage) GetRootUserPath() string {
	return s.rootUserPath
}

func (s *AccountsStorage) GetUserID() string {
	return s.userID
}

func (s *AccountsStorage) Save(account *Account) error {
	jsonBytes, err := json.MarshalIndent(account, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(s.accountFilePath, jsonBytes, filePerm)
}

func (s *AccountsStorage) LoadAccount(privateKey crypto.PrivateKey) *Account {
	fileBytes, err := os.ReadFile(s.accountFilePath)
	if err != nil {
		log.Errorf("Could not load file for account %s: %v\n", s.userID, err)
	}

	var account Account
	err = json.Unmarshal(fileBytes, &account)
	if err != nil {
		log.Errorf("Could not parse file for account %s: %v\n", s.userID, err)
	}

	account.key = privateKey

	if account.Registration == nil || account.Registration.Body.Status == "" {
		reg, err := tryRecoverRegistration(privateKey)
		if err != nil {
			log.Errorf("Could not load account for %s. Registration is nil: %#v\n", s.userID, err)
		}

		account.Registration = reg
		err = s.Save(&account)
		if err != nil {
			log.Errorf("Could not save account for %s. Registration is nil: %#v\n", s.userID, err)
		}
	}

	return &account
}

func (s *AccountsStorage) GetPrivateKey(keyType certcrypto.KeyType) (crypto.PrivateKey, error) {
	accKeyPath := filepath.Join(s.keysPath, s.userID+".key")

	if _, err := os.Stat(accKeyPath); os.IsNotExist(err) {
		//fmt.Printf("No key found for account %s. Generating a %s key.\n", s.userID, keyType)
		s.createKeysFolder()

		privateKey, err := generatePrivateKey(accKeyPath, keyType)
		if err != nil {
			return nil, fmt.Errorf(
				"could not generate RSA private account key for account %s: %v", s.userID, err,
			)
		}

		//fmt.Printf("Saved key to %s\n", accKeyPath)
		return privateKey, nil
	}

	privateKey, err := loadPrivateKey(accKeyPath)
	if err != nil {
		return nil, fmt.Errorf(
			fmt.Sprintf("could not load RSA private key from file %s: %v", accKeyPath, err),
		)
	}

	return privateKey, nil
}

func (s *AccountsStorage) createKeysFolder() {
	if err := createNonExistingFolder(s.keysPath); err != nil {
		log.Errorf("Could not check/create directory for account %s: %v\n", s.userID, err)
	}
}

func generatePrivateKey(file string, keyType certcrypto.KeyType) (crypto.PrivateKey, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(keyType)
	if err != nil {
		return nil, err
	}

	certOut, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	defer certOut.Close()

	pemKey := certcrypto.PEMBlock(privateKey)
	err = pem.Encode(certOut, pemKey)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadPrivateKey(file string) (crypto.PrivateKey, error) {
	keyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyBytes)

	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(keyBlock.Bytes)
	}

	return nil, errors.New("unknown private key type")
}

func tryRecoverRegistration(privateKey crypto.PrivateKey) (*registration.Resource, error) {
	acmeHost := acmeConfig.GetConfig().GetString("acme.host")
	if len(acmeHost) == 0 {
		acmeHost = acmeServer
	}

	acmeUserAgent := acmeConfig.GetConfig().GetString("acme.userAgent")
	if len(acmeUserAgent) == 0 {
		acmeUserAgent = "lego-cli/chargeover"
	}

	// couldn't load account but got a key. Try to look the account up.
	config := lego.NewConfig(&Account{key: privateKey})
	config.CADirURL = acmeHost
	config.UserAgent = acmeUserAgent

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	reg, err := client.Registration.ResolveAccountByKey()
	if err != nil {
		return nil, err
	}
	return reg, nil
}
