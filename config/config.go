package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) *viper.Viper {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatalf("error parsing configuration file: %v", err.Error())
	}

	return config
}

func GetConfig() *viper.Viper {
	return config
}
