package migration

import (
	"github.com/Celtech/ACME/web/database"
	"github.com/Celtech/ACME/web/model"
	log "github.com/sirupsen/logrus"
)

func RunMigrations() {
	err := database.GetDB().AutoMigrate(
		&model.Request{},
		&model.User{},
	)

	if err != nil {
		log.Errorf("There was a problem executing database migrations: %v", err)
	}
}
