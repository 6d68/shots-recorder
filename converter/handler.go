package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"os"
	"path"
)

type Handler func(context.Context, events.S3Event) error

func NewHandler(mp4Converter AviToMp4Converter, downUploader S3DownUploader) Handler {
	return func(ctx context.Context, event events.S3Event) error {

		record := event.Records[0]
		srcKey := record.S3.Object.Key
		srcBucket := record.S3.Bucket.Name

		// download
		sourceFile := path.Join(os.TempDir(), path.Base(srcKey))
		srcFile, err := os.Create(sourceFile)
		if err != nil {
			return err
		}

		err = downUploader.Download(srcFile, srcBucket, srcKey)
		if err != nil {
			return err
		}

		// convert
		convertedFile, err := mp4Converter.ToMp4(sourceFile)
		if err != nil {
			return err
		}

		// upload
		mp4File, err := os.Open(convertedFile)
		if err != nil {
			return err
		}

		targetBucket := os.Getenv("SHOTS_BUCKET")

		// keep original directory structure in target bucket
		uploadKey := path.Join(path.Dir(srcKey), path.Base(mp4File.Name()))
		err = downUploader.Upload(mp4File, targetBucket, uploadKey)

		if err != nil {
			return err
		}

		return nil
	}
}
