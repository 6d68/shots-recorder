package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"time"
)

type ShotsResponse struct {
	Key           string    `json:"key"`
	RecordingDate time.Time `json:"recording_date"`
}

type Response events.APIGatewayProxyResponse

type Handler func(context.Context) (Response, error)

func NewHandler() Handler {
	return func(ctx context.Context) (Response, error) {
		body, _ := json.Marshal([]ShotsResponse{
			{
				Key:           "shot-1",
				RecordingDate: time.Date(1999, 12, 21, 0, 0, 0, 0, time.Local),
			}, {
				Key:           "shot-2",
				RecordingDate: time.Date(1999, 12, 21, 0, 0, 0, 0, time.Local),
			},
		})

		response := Response{
			StatusCode:      200,
			IsBase64Encoded: false,
			Body:            string(body),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return response, nil
	}
}
