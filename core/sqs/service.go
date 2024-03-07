package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Config represents the configuration
type Config struct {
	Region string
}

// New initializes SNS service with default config
func New(cfg Config) *Service {
	s, err := session.NewSession(&aws.Config{Region: aws.String(cfg.Region)})
	if err != nil {
		panic(err)
	}

	return &Service{
		sqs: sqs.New(s),
		cfg: cfg,
	}
}

// Service represents the snsutil service
type Service struct {
	sqs *sqs.SQS
	cfg Config
}

// SQSMessageResponseData model
type SQSMessageResponse struct {
	Message          string `json:"Message"`
	MessageID        string `json:"MessageId"`
	Signature        string `json:"Signature"`
	SignatureVersion string `json:"SignatureVersion"`
	SigningCertURL   string `json:"SigningCertURL"`
	Timestamp        string `json:"Timestamp"`
	TopicArn         string `json:"TopicArn"`
	Type             string `json:"Type"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}
