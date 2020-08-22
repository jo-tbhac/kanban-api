package config

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

// ConfigList contains application information.
type ConfigList struct {
	AWS struct {
		Bucket string
		Region string
	}
	Database struct {
		User     string
		Name     string
		Host     string
		Password string
		Driver   string
		Log      bool
	}
	Web struct {
		Port   int
		Origin string
	}
}

var (
	// Config is instance of ConfigList.
	Config ConfigList
	sess   *session.Session
	err    error
)

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

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(Config.AWS.Region)},
	)

	if err != nil {
		log.Fatalf("invalid credentials: %v", err)
	}
}

// AWSSession returns a session of aws.
func AWSSession() *session.Session {
	return sess
}
