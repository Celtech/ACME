package main

import (
	"baker-acme/cmd"
	"baker-acme/config"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var conf *viper.Viper

func init() {
	var appEnv string
	if appEnv = os.Getenv("ACME_ENV"); appEnv == "" {
		appEnv = "development"
	}

	conf = config.Init(appEnv)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: !conf.GetBool("services.logger.color"),
		FullTimestamp: true,
	})
}

// @title           ACME API
// @version         1.0
// @description     Issue dynamic SSL certificates for 1+n domains.
// @description.markdown
// @termsOfService  https://chargeover.com/terms-of-service

// @contact.name   API Support
// @contact.url    https://git.rykelabs.com/rykelabs/acme-server/-/issues/new
// @contact.email  tim@chargeover.com

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
