package aws_modules_search

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            request.Body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}
