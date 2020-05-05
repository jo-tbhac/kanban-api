package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type ConfigList struct {
	Database struct {
		Name   string
		Driver string
	}
	Web struct {
		Port int
	}
}

var Config ConfigList

func init() {
	viper.SetConfigName("config")

	viper.SetConfigType("yml")

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("failed read config file: %v", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Printf("failed unmarshal config file: %v", err)
		os.Exit(1)
	}
}
