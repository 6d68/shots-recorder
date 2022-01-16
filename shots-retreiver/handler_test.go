package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {

	tests := []struct {
		name             string
		wantErr          assert.ErrorAssertionFunc
		wantResponseBody []ShotsResponse
		wantStatus       int
	}{
		{
			name:       "Should return two shots shot-1 and shot-2",
			wantErr:    assert.NoError,
			wantStatus: http.StatusOK,
			wantResponseBody: []ShotsResponse{
				{
					Key:           "shot-1",
					RecordingDate: time.Date(1999, 12, 21, 0, 0, 0, 0, time.Local),
				}, {
					Key:           "shot-2",
					RecordingDate: time.Date(1999, 12, 21, 0, 0, 0, 0, time.Local),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			handler := NewHandler()
			response, err := handler(nil)

			var shots []ShotsResponse
			json.Unmarshal([]byte(response.Body), &shots)

			assert.EqualValues(t, tt.wantResponseBody, shots)
			tt.wantErr(t, err)
		})
	}
}
