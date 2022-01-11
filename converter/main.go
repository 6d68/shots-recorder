package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	manager := &FileManager{
		D: s3manager.NewDownloader(sess),
		U: s3manager.NewUploader(sess),
	}

	converter := Converter{FfMpegPath: "/opt/bin/ffmpeg"}
	handler := NewHandler(converter, manager)

	lambda.Start(handler)
}
