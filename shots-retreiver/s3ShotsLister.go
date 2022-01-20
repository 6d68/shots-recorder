package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"os"
	"strings"
)

type shotsLister interface {
	All() ([]ShotsMeta, error)
}

type s3ShotsLister struct {
	s3 s3iface.S3API
}

type ShotsMeta struct {
	key string
}

func (s s3ShotsLister) All() ([]ShotsMeta, error) {
	bucket := os.Getenv("SHOTS_BUCKET")
	listResp, err := s.s3.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &bucket})
	if err != nil {
		return nil, fmt.Errorf("error listing objects of bucket %s: %w", bucket, err)
	}

	var metas []ShotsMeta
	for _, item := range listResp.Contents {
		key := fmt.Sprint(*item.Key)

		if strings.HasSuffix(key, ".mp4") {
			metas = append(metas, ShotsMeta{key: key})
		}
	}

	return metas, nil
}
