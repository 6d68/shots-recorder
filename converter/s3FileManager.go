package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"io"
)

type S3DownUploader interface {
	Download(writer io.WriterAt, bucket string, key string) error
	Upload(reader io.Reader, bucket string, key string) error
}

type FileManager struct {
	D s3manageriface.DownloaderAPI
	U s3manageriface.UploaderAPI
}

func (f FileManager) Download(writer io.WriterAt, bucket string, fileKey string) error {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}

	_, err := f.D.Download(writer, getObjectInput)

	if err != nil {
		return fmt.Errorf("error downloading writer with key %s from bucket %s: %w", fileKey, bucket, err)
	}

	return nil
}

func (f FileManager) Upload(reader io.Reader, bucket string, fileKey string) error {
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   reader,
	}

	_, err := f.U.Upload(uploadInput)

	if err != nil {
		return fmt.Errorf("error uploading file with key %s to %s, %w", fileKey, bucket, err)
	}

	return nil
}
