package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"testing"
)

type mockS3 struct {
	s3manageriface.DownloaderAPI
	s3manageriface.UploaderAPI
	mock.Mock
}

func (m *mockS3) Download(file io.WriterAt, input *s3.GetObjectInput, _ ...func(*s3manager.Downloader)) (int64, error) {
	args := m.Called(file, input)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockS3) Upload(input *s3manager.UploadInput, _ ...func(uploader *s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3manager.UploadOutput), args.Error(1)
}

func setup() (*mockS3, *FileManager) {
	s3Mock := new(mockS3)
	mockFileManager := &FileManager{
		D: s3Mock,
		U: s3Mock,
	}
	return s3Mock, mockFileManager
}

func TestFileManager_Download(t *testing.T) {

	type args struct {
		bucket  string
		fileKey string
	}
	tests := []struct {
		name         string
		args         args
		mockInput    *s3.GetObjectInput
		mockErrorOut error
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			args: args{
				bucket:  "bucket",
				fileKey: "key",
			},
			mockInput: &s3.GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			mockErrorOut: nil,
			wantErr:      assert.NoError,
		}, {
			args: args{
				bucket:  "bucket",
				fileKey: "key",
			},
			mockInput: &s3.GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			mockErrorOut: errors.New("failed to download"),
			wantErr:      assert.Error,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			s3Manager, f := setup()

			s3Manager.On("Download", aws.NewWriteAtBuffer([]byte{}), tt.mockInput).Return(int64(0), tt.mockErrorOut).Once()

			tt.wantErr(t, f.Download(aws.NewWriteAtBuffer([]byte{}), tt.args.bucket, tt.args.fileKey), fmt.Sprintf("Download(%v, %v, %v)", aws.NewWriteAtBuffer([]byte{}), tt.args.bucket, tt.args.fileKey))
		})
	}
}

func TestFileManager_Upload(t *testing.T) {

	type args struct {
		bucket  string
		fileKey string
	}

	tests := []struct {
		name         string
		args         args
		input        *s3manager.UploadInput
		mockErrorOut error
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			args: args{
				bucket:  "bucket",
				fileKey: "key",
			},
			input: &s3manager.UploadInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   bytes.NewReader([]byte{}),
			},
			mockErrorOut: nil,
			wantErr:      assert.NoError,
		}, {
			args: args{
				bucket:  "bucket",
				fileKey: "key",
			},
			input: &s3manager.UploadInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   bytes.NewReader([]byte{}),
			},
			mockErrorOut: errors.New("failed to download"),
			wantErr:      assert.Error,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			s3Manager, f := setup()

			s3Manager.On("Upload", tt.input).Return(&s3manager.UploadOutput{}, tt.mockErrorOut).Once()

			tt.wantErr(t, f.Upload(bytes.NewReader([]byte{}), tt.args.bucket, tt.args.fileKey), fmt.Sprintf("Download(%v, %v, %v)", bytes.NewReader([]byte{}), tt.args.bucket, tt.args.fileKey))
		})
	}
}
