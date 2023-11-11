package aws_discovery

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"net/http"
	"terraform-serverless-private-registry/lib/helpers"
)

var (
	logger *zap.Logger
)

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
}

type Response struct {
	Modules   string `json:"modules.v1"`
	Providers string `json:"providers.v1"`
	Custom    string `json:"custom.v1"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()
	logger.Debug("discovery called",
		zap.String("reqId", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	resp := Response{
		Modules:   "/v1/modules/",
		Providers: "/v1/providers/",
		Custom:    "/v1/custom/",
	}

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
