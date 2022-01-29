package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"os"
	"time"
)

type UrlSigner interface {
	Presign(key string) (string, error)
}

type S3UrlPreSigner struct {
	s3 s3iface.S3API
}

func (p S3UrlPreSigner) Presign(key string) (string, error) {
	req, _ := p.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("SHOTS_BUCKET")),
		Key:    aws.String(key),
	})

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("error signing object access url: %w", err)
	}

	return url, nil
}
