package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"testing"
)

type mockConverter struct {
	AviToMp4Converter
	mockErr error
}

func (m *mockConverter) ToMp4(in string) (string, error) {

	if m.mockErr != nil {
		return "", m.mockErr
	}

	mp4File := setExtensions(in, ".mp4")
	f, err := os.Create(mp4File)
	defer f.Close()
	if err != nil {
		return "", err
	}

	return mp4File, nil
}

func setupMockConverter(err error) *mockConverter {
	return &mockConverter{
		mockErr: err,
	}
}

type mockDownUploader struct {
	S3DownUploader
	mock.Mock
}

func (m mockDownUploader) Download(writer io.WriterAt, bucket string, key string) error {
	args := m.Called(writer, bucket, key)
	return args.Error(0)
}

func (m mockDownUploader) Upload(reader io.Reader, bucket string, key string) error {
	args := m.Called(reader, bucket, key)
	return args.Error(0)
}

func setupMockDownUploader() *mockDownUploader {
	return new(mockDownUploader)
}

func TestNewHandler(t *testing.T) {
	type args struct {
		event events.S3Event
	}

	tests := []struct {
		name             string
		args             args
		wantErr          assert.ErrorAssertionFunc
		wantUploadBucket string
		wantUploadKey    string
		downloadErr      error
		uploadErr        error
		conversionErr    error
	}{
		{
			name:             "Should convert newly uploaded avi to mp4 and upload to shots_bucket",
			args:             args{event: s3Event("bucket", "cam1/20220108/recording.avi")},
			downloadErr:      nil,
			wantUploadBucket: "shots_bucket",
			wantUploadKey:    "cam1/20220108/recording.mp4",
			wantErr:          assert.NoError,
		}, {
			name:             "Should abort conversion when s3 download of avi not successful",
			args:             args{event: s3Event("bucket", "cam1/20220108/recording.avi")},
			downloadErr:      errors.New("download failed"),
			wantUploadBucket: "shots_bucket",
			wantUploadKey:    "cam1/20220108/recording.mp4",
			wantErr:          assert.Error,
		}, {
			name:             "Should abort conversion when s3 upload of avi not successful",
			args:             args{event: s3Event("bucket", "cam1/20220108/recording.avi")},
			uploadErr:        errors.New("upload failed"),
			wantUploadBucket: "shots_bucket",
			wantUploadKey:    "cam1/20220108/recording.mp4",
			wantErr:          assert.Error,
		}, {
			name:             "Should abort when conversion fails",
			args:             args{event: s3Event("bucket", "cam1/20220108/recording.avi")},
			conversionErr:    errors.New("conversion failed"),
			wantUploadBucket: "shots_bucket",
			wantUploadKey:    "cam1/20220108/recording.mp4",
			wantErr:          assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			converterMock := setupMockConverter(tt.conversionErr)
			downUploader := setupMockDownUploader()
			os.Setenv("SHOTS_BUCKET", tt.wantUploadBucket)
			bucket, key := firstEventRecordBucketAndKey(tt.args.event)

			downUploader.On("Download", mock.Anything, bucket, key).Return(tt.downloadErr)
			downUploader.On("Upload", mock.Anything, tt.wantUploadBucket, tt.wantUploadKey).Return(tt.uploadErr)

			s3EventHandler := NewHandler(converterMock, downUploader)
			tt.wantErr(t, s3EventHandler(nil, tt.args.event))
		})
	}
}

func s3Event(bucket string, key string) events.S3Event {
	return events.S3Event{Records: []events.S3EventRecord{
		{
			S3: events.S3Entity{
				SchemaVersion:   "",
				ConfigurationID: "",
				Bucket: events.S3Bucket{
					Name:          bucket,
					OwnerIdentity: events.S3UserIdentity{},
					Arn:           "",
				},
				Object: events.S3Object{
					Key: key,
				},
			},
		},
	}}
}

func firstEventRecordBucketAndKey(event events.S3Event) (string, string) {
	r := event.Records[0]
	return r.S3.Bucket.Name, r.S3.Object.Key
}

func Test_setExtensions(t *testing.T) {
	type args struct {
		fileName string
		newExt   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should change extension form txt to mp4",
			args: args{
				fileName: "test.txt",
				newExt:   ".mp4",
			},
			want: "test.mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, setExtensions(tt.args.fileName, tt.args.newExt), "setExtensions(%v, %v)", tt.args.fileName, tt.args.newExt)
		})
	}
}
