package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type ShotsResponse struct {
	Key string `json:"key"`
}

type Response events.APIGatewayProxyResponse

type ErrorDetails struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type Handler func(context.Context) (Response, error)

func NewHandler(shotsLister shotsLister) Handler {
	return func(ctx context.Context) (Response, error) {
		metas, err := shotsLister.All()
		if err != nil {
			log.Println(fmt.Errorf("error getting shots meta from storage: %w", err))
			return internalServerError(), err
		}

		var shots []ShotsResponse

		for _, meta := range metas {
			shots = append(shots, ShotsResponse{Key: meta.key})
		}

		body, err := json.Marshal(shots)
		if err != nil {
			log.Println(fmt.Errorf("error json.Marshal(shots): %w", err))
			return internalServerError(), err
		}

		return Response{
			StatusCode:      200,
			IsBase64Encoded: false,
			Body:            string(body),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}

func internalServerError() Response {
	status := http.StatusInternalServerError
	response := ErrorResponse{Error: ErrorDetails{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
	}}

	body, _ := json.Marshal(response)

	return Response{
		StatusCode:      status,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
