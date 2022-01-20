package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func main() {

	region := os.Getenv("REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		log.Fatalf("Error initializing AWS session: %v", err)
	}

	lister := s3ShotsLister{s3: s3.New(sess)}

	handler := NewHandler(lister)
	lambda.Start(handler)
}
