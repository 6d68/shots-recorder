package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"strings"
)

type ShotsSignRequest struct {
	Key string `json:"key"`
}

type ShotsSignResponse struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

type Response events.APIGatewayProxyResponse

type Request events.APIGatewayProxyRequest

type ErrorDetails struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type Handler func(context.Context, Request) (Response, error)

func NewHandler(preSigner UrlSigner) Handler {
	return func(ctx context.Context, request Request) (Response, error) {

		signRequest := ShotsSignRequest{}

		if err := json.Unmarshal([]byte(request.Body), &signRequest); err != nil {
			return errorResponse(http.StatusBadRequest, "provided json body not valid"), err
		}

		if strings.TrimSpace(signRequest.Key) == "" {
			return errorResponse(http.StatusBadRequest, "signRequest.key must not be empty"), errors.New("signRequest.key must not be empty")
		}

		url, err := preSigner.
			Presign(signRequest.Key)
		if err != nil {
			return errorResponse(http.StatusInternalServerError, "internal server error"), err
		}

		body, _ := json.Marshal(ShotsSignResponse{Key: signRequest.Key, Url: url})

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

func errorResponse(status int, message string) Response {
	response := ErrorResponse{Error: ErrorDetails{
		Status:  status,
		Message: message,
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
