package database

import (
	"github.com/Celtech/ACME/config"
	"github.com/glebarez/sqlite"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() *gorm.DB {
	basePath := config.GetConfig().GetString("acme.dataPath")
	if len(basePath) == 0 {
		basePath = "/data"
	}

	dbFilePath := filepath.Join(basePath, "ssl-certify.db")
	d, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}
	db = d

	return db
}

func GetDB() *gorm.DB {
	return db
}
