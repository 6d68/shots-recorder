package main

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type mockUrlSigner struct {
	S3UrlPreSigner
	mock.Mock
}

func (m *mockUrlSigner) Presign(url string) (string, error) {
	args := m.Called(url)
	return args.Get(0).(string), args.Error(1)
}

func TestNewHandler(t *testing.T) {

	tests := []struct {
		name                 string
		key                  string
		shotsSignRequestBody string
		presignError         error
		wantStatus           int
		wantErr              require.ErrorAssertionFunc
		wantResponseBody     ShotsSignResponse
		isMalformedRequest   bool
	}{
		{
			name:                 "Should return return internal server error when presigning fails",
			key:                  "cam-1/20220113/112349_1126445.mp4",
			shotsSignRequestBody: "{\n\t\"key\": \"cam-1/20220113/112349_1126445.mp4\"\n}",
			presignError:         errors.New("some error"),
			wantStatus:           http.StatusInternalServerError,
			wantResponseBody:     ShotsSignResponse{},
			wantErr:              require.Error,
		},
		{
			name:                 "Should create temp url if presigning was successful",
			key:                  "cam-1/20220113/112349_1126445.mp4",
			shotsSignRequestBody: "{\n\t\"key\": \"cam-1/20220113/112349_1126445.mp4\"\n}",
			presignError:         nil,
			wantStatus:           http.StatusOK,
			wantResponseBody: ShotsSignResponse{
				Key: "cam-1/20220113/112349_1126445.mp4",
				Url: "cam-1/20220113/112349_1126445.mp4",
			},
			wantErr: require.NoError,
		},
		{
			name:                 "Should return return bad request when input is not valid",
			key:                  "cam-1/20220113/112349_1126445.mp4",
			shotsSignRequestBody: "{\n\t\"bad_idea\": \"cam-1/20220113/112349_1126445.mp4\"\n}",
			presignError:         nil,
			wantStatus:           http.StatusBadRequest,
			wantResponseBody:     ShotsSignResponse{},
			wantErr:              require.Error,
		},
		{
			name:                 "Should return return bad request when json at all is not valid",
			key:                  "cam-1/20220113/112349_1126445.mp4",
			shotsSignRequestBody: "{}}",
			presignError:         nil,
			wantStatus:           http.StatusBadRequest,
			wantResponseBody:     ShotsSignResponse{},
			wantErr:              require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			signerMock := new(mockUrlSigner)
			signerMock.On("Presign", tt.key).Return(tt.key, tt.presignError)

			handler := NewHandler(signerMock)

			// Act
			request := Request{
				Body: tt.shotsSignRequestBody,
			}
			response, err := handler(nil, request)

			// Assert
			tt.wantErr(t, err)

			require.Equal(t, tt.wantStatus, response.StatusCode)

			if tt.wantStatus == 200 {
				var responseBody ShotsSignResponse
				err := json.Unmarshal([]byte(response.Body), &responseBody)

				require.NoError(t, err)
				require.EqualValues(t, tt.wantResponseBody, responseBody)
			}
		})
	}
}
