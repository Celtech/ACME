package acme

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/Celtech/ACME/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Celtech/ACME/internal/acme/certificate"
	certcrypto "github.com/Celtech/ACME/internal/acme/crypto"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/idna"
)

const (
	baseCertificatesFolderName = "certificates"
	baseArchivesFolderName     = "archives"
)

// CertificatesStorage a certificates storage.
//
// rootPath:
//
//	./.lego/certificates/
//	     │      └── root certificates directory
//	     └── "path" option
//
// archivePath:
//
//	./.lego/archives/
//	     │      └── archived certificates directory
//	     └── "path" option
type CertificatesStorage struct {
	rootPath    string
	archivePath string
	pem         bool
	pfx         bool
	pfxPassword string
}

// NewCertificatesStorage create a new certificates storage.
func NewCertificatesStorage() *CertificatesStorage {
	basePath := config.GetConfig().GetString("acme.dataPath")
	if len(basePath) == 0 {
		basePath = "/data"
	}

	return &CertificatesStorage{
		rootPath:    filepath.Join(basePath, baseCertificatesFolderName),
		archivePath: filepath.Join(basePath, baseArchivesFolderName),
		pem:         true,
		pfx:         false,
		pfxPassword: "",
	}
}

func (s *CertificatesStorage) CreateRootFolder() error {
	return createNonExistingFolder(s.rootPath)
}

func (s *CertificatesStorage) CreateArchiveFolder() {
	err := createNonExistingFolder(s.archivePath)
	if err != nil {
		log.Errorf("Could not check/create path: %v", err)
	}
}

func (s *CertificatesStorage) GetRootPath() string {
	return s.rootPath
}

func (s *CertificatesStorage) SaveResource(certRes *certificate.Resource) error {
	domain := certRes.Domain

	// We store the certificate, private key and metadata in different files
	// as web servers would not be able to work with a combined file.
	err := s.WriteFile(domain, ".crt", certRes.Certificate)
	if err != nil {
		return fmt.Errorf(
			"unable to save certificate for domain %s\n\t%v", domain, err,
		)
	}

	if certRes.IssuerCertificate != nil {
		err = s.WriteFile(domain, ".issuer.crt", certRes.IssuerCertificate)
		if err != nil {
			return fmt.Errorf(
				"unable to save IssuerCertificate for domain %s\n\t%v", domain, err,
			)
		}
	}

	// if we were given a CSR, we don't know the private key
	if certRes.PrivateKey != nil {
		err = s.WriteCertificateFiles(domain, certRes)
		if err != nil {
			return fmt.Errorf(
				"unable to save PrivateKey for domain %s\n\t%v", domain, err,
			)
		}
	} else if s.pem || s.pfx {
		return fmt.Errorf(
			"unable to save PEM without private key for domain %s. Are you using a CSR?", domain,
		)
	}

	jsonBytes, err := json.MarshalIndent(certRes, "", "\t")
	if err != nil {
		return fmt.Errorf(
			"unable to marshal CertResource for domain %s\n\t%v", domain, err,
		)
	}

	err = s.WriteFile(domain, ".json", jsonBytes)
	if err != nil {
		return fmt.Errorf(
			"unable to save CertResource for domain %s\n\t%v", domain, err,
		)
	}

	return nil
}

func (s *CertificatesStorage) ReadResource(domain string) certificate.Resource {
	raw, err := s.ReadFile(domain, ".json")
	if err != nil {
		log.Errorf("error while loading the meta data for domain %s\n\t%v", domain, err)
	}

	var resource certificate.Resource
	if err = json.Unmarshal(raw, &resource); err != nil {
		log.Errorf("error while marshaling the meta data for domain %s\n\t%v", domain, err)
	}

	return resource
}

func (s *CertificatesStorage) ExistsFile(domain, extension string) bool {
	filePath := s.GetFileName(domain, extension)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Error(err)
	}
	return true
}

func (s *CertificatesStorage) ReadFile(domain, extension string) ([]byte, error) {
	return os.ReadFile(s.GetFileName(domain, extension))
}

func (s *CertificatesStorage) GetFileName(domain, extension string) string {
	filename := sanitizedDomain(domain) + extension
	return filepath.Join(s.rootPath, filename)
}

func (s *CertificatesStorage) ReadCertificate(domain, extension string) ([]*x509.Certificate, error) {
	content, err := s.ReadFile(domain, extension)
	if err != nil {
		return nil, err
	}

	// The input may be a bundle or a single certificate.
	return certcrypto.ParsePEMBundle(content)
}

// CheckCrtList will check to see if the cert list file contains our given file or not. It will
// return an (int greater than 0 if found) or (0 on an error or if it wasn't found)
func (s *CertificatesStorage) CheckCrtList(certFileName string, crtList *os.File) (int, error) {
	scanner := bufio.NewScanner(crtList)
	line := 1

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), certFileName) {
			return line, nil
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, nil
}

func (s *CertificatesStorage) GetCrtListPath() string {
	return filepath.Join(s.rootPath, "crt-list.txt")
}

// WriteCrtList will write a newly issued cert to the crt-list.txt file to be used by reverse proxies such as HAProxy
// Note this will also check to see if the file already exists, if it does, it won't write anything to prevent
// duplication
func (s *CertificatesStorage) WriteCrtList(domain string) error {
	filePath := s.GetCrtListPath()
	crtList, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer crtList.Close()

	var baseFileName = sanitizedDomain(domain)
	var certFileName = filepath.Join(s.rootPath, baseFileName+".pem\n")
	lineNo, err := s.CheckCrtList(certFileName, crtList)

	if lineNo == 0 && err == nil {
		_, writeErr := crtList.WriteString(certFileName)
		return writeErr
	}

	return err
}

func (s *CertificatesStorage) WriteFile(domain, extension string, data []byte) error {
	var baseFileName = sanitizedDomain(domain)
	filePath := filepath.Join(s.rootPath, baseFileName+extension)

	return os.WriteFile(filePath, data, filePerm)
}

func (s *CertificatesStorage) WriteCertificateFiles(domain string, certRes *certificate.Resource) error {
	err := s.WriteFile(domain, ".key", certRes.PrivateKey)
	if err != nil {
		return fmt.Errorf("unable to save key file: %w", err)
	}

	if s.pem {
		err = s.WriteFile(domain, ".pem", bytes.Join([][]byte{certRes.Certificate, certRes.PrivateKey}, nil))
		if err != nil {
			return fmt.Errorf("unable to save PEM file: %w", err)
		}

		err = s.WriteCrtList(domain)
		if err != nil {
			return fmt.Errorf("unable to write certificate to crt-list.txt: %w", err)
		}
	}

	return nil
}

func (s *CertificatesStorage) MoveToArchive(domain string) error {
	matches, err := filepath.Glob(filepath.Join(s.rootPath, sanitizedDomain(domain)+".*"))
	if err != nil {
		return err
	}

	for _, oldFile := range matches {
		date := strconv.FormatInt(time.Now().Unix(), 10)
		filename := date + "." + filepath.Base(oldFile)
		newFile := filepath.Join(s.archivePath, filename)

		err = os.Rename(oldFile, newFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// sanitizedDomain Make sure no funny chars are in the cert names (like wildcards ;)).
func sanitizedDomain(domain string) string {
	safe, err := idna.ToASCII(strings.ReplaceAll(domain, "*", "_"))
	if err != nil {
		log.Error(err)
	}
	return safe
}
