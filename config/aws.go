package config

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	sess *session.Session
	err  error
)

func init() {
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)

	if err != nil {
		log.Fatalf("invalid credentials: %v", err)
	}
}

// AWSSession returns a session of aws.
func AWSSession() *session.Session {
	return sess
}
