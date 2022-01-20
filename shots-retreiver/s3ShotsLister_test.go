package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type s3Mock struct {
	s3iface.S3API
	mock.Mock
}

func (m *s3Mock) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func setupShotsListerMock() (*s3Mock, shotsLister) {
	mockS3 := new(s3Mock)
	lister := s3ShotsLister{s3: mockS3}

	return mockS3, lister
}

func Test_s3ShotsLister_All(t *testing.T) {

	tests := []struct {
		name      string
		want      []ShotsMeta
		wantErr   assert.ErrorAssertionFunc
		s3Content []*s3.Object
		s3Error   error
	}{
		{
			name:      "should only return mp4 video files",
			s3Content: []*s3.Object{{Key: aws.String("key.mp4")}, {Key: aws.String("some key")}},
			s3Error:   nil,
			want:      []ShotsMeta{{key: "key.mp4"}},
			wantErr:   assert.NoError,
		},
		{
			name:      "should return error only when s3 access fails",
			s3Content: nil,
			s3Error:   errors.New("some error"),
			want:      nil,
			wantErr:   assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockS3, s := setupShotsListerMock()

			mockListObjectsV2Output := &s3.ListObjectsV2Output{
				Contents: tt.s3Content,
			}

			mockS3.On("ListObjectsV2", &s3.ListObjectsV2Input{Bucket: aws.String("")}).Return(mockListObjectsV2Output, tt.s3Error)
			got, err := s.All()
			if !tt.wantErr(t, err, fmt.Sprintf("All()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "All()")
		})
	}
}
