package acme

import (
	"context"
	"errors"
	"fmt"
	"github.com/Celtech/ACME/internal/acme/plugins"
)

func RunPlugins(ctx context.Context) error {
	if ctx.Value("domain") == nil {
		return fmt.Errorf("`domain` must be set on the RunPlugins context")
	}

	certsStorage := NewCertificatesStorage()
	domain := ctx.Value("domain").(string)

	if !certsStorage.ExistsFile(domain, ".pem") {
		return errors.New("certificate does not exist")
	}

	path := certsStorage.GetFileName(domain, ".pem")
	certListPath := certsStorage.GetCrtListPath()
	contents, err := certsStorage.ReadFile(domain, ".pem")
	if err != nil {
		return err
	}

	return plugins.Run(certListPath, path, string(contents))
}
