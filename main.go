package main

import (
	"github.com/Celtech/ACME/cmd"
	"github.com/Celtech/ACME/config"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	var appEnv string
	if appEnv = os.Getenv("ACME_ENV"); appEnv == "" {
		appEnv = "development"
	}

	conf := config.Init(appEnv)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: !conf.GetBool("services.logger.color"),
		FullTimestamp: true,
		ForceColors:   conf.GetBool("services.logger.color"),
	})
}

// @title           ACME API
// @version         1.0
// @description     Issue dynamic SSL certificates for 1+n domains.
// @description.markdown

// @license.name  MIT
// @license.url   https://git.rykelabs.com/rykelabs/acme-server/-/blob/main/LICENSE.md

// @host      acme.chargeover.com:9022
// @BasePath  /api/v1

// @accept	json
// @produce json

// @schemes https

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description					JWT Authorization Token

func main() {
	cmd.Execute()
}
