package aws_index

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"terraform-serverless-private-registry/lib/helpers"
)

type Response struct {
	Status string `json:"status"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	resp := new(Response)
	resp.Status = "OK"
	return helpers.ApiResponse(http.StatusOK, resp), nil
}
