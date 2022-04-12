package http01

import (
	"fmt"

	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/acme/api"
	"github.com/go-acme/lego/v4/challenge"
	log "github.com/sirupsen/logrus"
)

type ValidateFunc func(core *api.Core, domain string, chlng acme.Challenge) error

// ChallengePath returns the URL path for the `http-01` challenge.
func ChallengePath(token string) string {
	return "/.well-known/acme-challenge/" + token
}

type Challenge struct {
	core     *api.Core
	validate ValidateFunc
	provider challenge.Provider
}

func NewChallenge(core *api.Core, validate ValidateFunc, provider challenge.Provider) *Challenge {
	return &Challenge{
		core:     core,
		validate: validate,
		provider: provider,
	}
}

func (c *Challenge) SetProvider(provider challenge.Provider) {
	c.provider = provider
}

func (c *Challenge) Solve(authz acme.Authorization) error {
	domain := challenge.GetTargetedDomain(authz)
	log.WithFields(log.Fields{
		"domain": domain,
	}).Infof("acme: Trying to solve HTTP-01")

	chlng, err := challenge.FindChallenge(challenge.HTTP01, authz)
	if err != nil {
		return err
	}

	// Generate the Key Authorization for the challenge
	keyAuth, err := c.core.GetKeyAuthorization(chlng.Token)
	if err != nil {
		return err
	}

	err = c.provider.Present(authz.Identifier.Value, chlng.Token, keyAuth)
	if err != nil {
		return fmt.Errorf("[%s] acme: error presenting token: %w", domain, err)
	}
	defer func() {
		err := c.provider.CleanUp(authz.Identifier.Value, chlng.Token, keyAuth)
		if err != nil {
			log.WithFields(log.Fields{
				"domain": domain,
			}).Warnf("acme: cleaning up failed: %v", err)
		}
	}()

	chlng.KeyAuthorization = keyAuth
	return c.validate(c.core, domain, chlng)
}
