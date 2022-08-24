package main

import (
	"github.com/Celtech/ACME/cmd"
	"github.com/Celtech/ACME/config"
	"github.com/Celtech/ACME/internal/queue"
	"github.com/Celtech/ACME/web/database"
	"github.com/Celtech/ACME/web/database/migration"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	Version string
	Date    string
)

func branding() {
	log.Infof(`
   _____ _____ _         _____          _   _  __
  / ____/ ____| |       / ____|        | | (_)/ _|
 | (___| (___ | |      | |     ___ _ __| |_ _| |_ _   _
  \___ \\___ \| |      | |    / _ \ '__| __| |  _| | | |
  ____) |___) | |____  | |___|  __/ |  | |_| | | | |_| |
 |_____/_____/|______|  \_____\___|_|   \__|_|_|  \__, |
                                                   __/ |
                                                  |___/
`)

	log.Infof("Starting SSL Certify version %s built on %s", Version, Date)
}

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

	branding()

	database.Init()
	migration.RunMigrations()
	queue.QueueMgr = queue.NewQueue(conf.GetString("redis.name"))
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
