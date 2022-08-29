package database

import (
	"fmt"
	"github.com/Celtech/ACME/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() *gorm.DB {
	conf := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.GetString("database.user"),
		conf.GetString("database.pass"),
		conf.GetString("database.host"),
		conf.GetString("database.port"),
		conf.GetString("database.name"),
	)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}
	db = d

	return db
}

func GetDB() *gorm.DB {
	return db
}
