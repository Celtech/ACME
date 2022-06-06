package config

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) *viper.Viper {
	curDir, e := os.Getwd()
	if e != nil {
		log.Error("there was a project getting the current working directory")
	}

	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath(fmt.Sprintf("%s/config", getConfigDir(curDir)))
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

func getConfigDir(dir string) string {
	if _, err := os.Stat(fmt.Sprintf("%s/main.go", dir)); errors.Is(err, os.ErrNotExist) {
		return getConfigDir(fmt.Sprintf("%s/..", dir))
	} else {
		return dir
	}
}
