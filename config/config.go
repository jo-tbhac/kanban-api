package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// ConfigList contains application information.
type ConfigList struct {
	AWS struct {
		Bucket string
	}
	Database struct {
		Name   string
		Driver string
	}
	Web struct {
		Port   int
		Origin string
	}
}

// Config is instance of ConfigList.
var Config ConfigList

func init() {
	if env := os.Getenv("environment"); env == "test" {
		return
	}

	viper.SetConfigName("config")

	viper.SetConfigType("yml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("failed read config file: %v", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Printf("failed unmarshal config file: %v", err)
		os.Exit(1)
	}
}
