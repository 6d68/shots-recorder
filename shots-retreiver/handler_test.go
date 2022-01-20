package main

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type mockShotsLister struct {
	shotsLister
	mock.Mock
}

func (m *mockShotsLister) All() ([]ShotsMeta, error) {
	args := m.Called()
	return args.Get(0).([]ShotsMeta), args.Error(1)
}

func TestNewHandler(t *testing.T) {

	tests := []struct {
		name       string
		shotsMeta  []ShotsMeta
		shotsError error
		wantStatus int
		wantShots  []ShotsResponse
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "Should return two shots shot-1 and shot-2",
			shotsMeta:  []ShotsMeta{{key: "shot-1"}, {key: "shot-2"}},
			shotsError: nil,
			wantStatus: http.StatusOK,
			wantShots: []ShotsResponse{
				{
					Key: "shot-1",
				}, {
					Key: "shot-2",
				},
			},
			wantErr: assert.NoError,
		}, {
			name:       "Should return return internal server error when reading shots fails",
			shotsMeta:  nil,
			shotsError: errors.New("some error"),
			wantStatus: http.StatusInternalServerError,
			wantShots:  nil,
			wantErr:    assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			shotsListerMock := new(mockShotsLister)
			shotsListerMock.On("All").Return(tt.shotsMeta, tt.shotsError)
			handler := NewHandler(shotsListerMock)

			// Act
			response, err := handler(nil)

			// Assert
			tt.wantErr(t, err)
			var shots []ShotsResponse
			json.Unmarshal([]byte(response.Body), &shots)
			assert.EqualValues(t, tt.wantShots, shots)
		})
	}
}
