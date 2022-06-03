package migration

import (
	"baker-acme/web/database"
	"baker-acme/web/model"
)

func RunMigrations() {
	database.GetDB().AutoMigrate(&model.Request{})
}
