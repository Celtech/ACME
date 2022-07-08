package migration

import (
	"github.com/Celtech/ACME/web/database"
	"github.com/Celtech/ACME/web/model"
)

func RunMigrations() {
	database.GetDB().AutoMigrate(
		&model.Request{},
		&model.User{},
	)
}
