package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"net/http"
	"os"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/modules"
	"terraform-serverless-private-registry/lib/storage"
)

var (
	modulesSvc *modules.Modules
	logger     *zap.Logger
)

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
	bucketName := os.Getenv("BUCKET_NAME")
	storage, _ := storage.NewStorage(bucketName, logger)
	modulesSvc, _ = modules.NewModules(storage, logger)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()
	logger.Debug(fmt.Sprintf("%s lambda called", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	namespace := request.PathParameters["namespace"]
	name := request.PathParameters["name"]
	provider := request.PathParameters["provider"]

	params := modules.ModuleParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
	}
	resp, err := modulesSvc.ListModuleVersions(request.RequestContext.RequestID, params)

	if err != nil {
		logger.Error(fmt.Sprintf("%s error %s", request.RequestContext.RequestID, err.Error()),
			zap.Error(err),
		)
		return helpers.ApiErrorUnknown(request.RequestContext.RequestID), nil
	}

	if resp.ModuleExists == false {
		return helpers.ApiErrorNotFound(request.RequestContext.RequestID, ""), nil
	}

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
