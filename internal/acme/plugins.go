package acme

import (
	"errors"
	"github.com/Celtech/ACME/internal/acme/plugins"
)

func RunPlugins(domain string, cert *CertificatesStorage) error {
	if !cert.ExistsFile(domain, ".pem") {
		return errors.New("certificate does not exist")
	}

	contents, err := cert.ReadFile(domain, ".pem")
	if err != nil {
		return err
	}

	path := cert.GetFileName(domain, ".pem")

	certListPath := cert.GetCrtListPath()

	return plugins.Run(certListPath, path, string(contents))
}
