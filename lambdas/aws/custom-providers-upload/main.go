package aws_custom_providers_upload

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"net/http"
	"os"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/providers"
	"terraform-serverless-private-registry/lib/storage"
)

var (
	providersSvc *providers.Providers
	logger       *zap.Logger
)

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
	bucketName := os.Getenv("BUCKET_NAME")
	storageSvc, _ := storage.NewStorage(bucketName, logger)
	providersSvc, _ = providers.NewProviders(storageSvc, logger)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()
	logger.Debug(fmt.Sprintf("%s lambda called", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	providerNamespace := request.PathParameters["namespace"]
	providerType := request.PathParameters["type"]
	providerVersion := request.PathParameters["version"]
	providerOs := request.PathParameters["os"]
	providerArch := request.PathParameters["arch"]

	params := providers.GetUploadInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
		Version:   &providerVersion,
		OS:        &providerOs,
		Arch:      &providerArch,
	}
	resp, err := providersSvc.GetUpload(request.RequestContext.RequestID, params)

	if err != nil {
		logger.Error(fmt.Sprintf("%s error %s", request.RequestContext.RequestID, err.Error()),
			zap.Error(err),
		)
		return helpers.ApiErrorUnknown(request.RequestContext.RequestID), nil
	}

	if resp.ProviderExists != nil && *resp.ProviderExists == true {
		return helpers.ApiErrorConflict(request.RequestContext.RequestID, ""), nil
	}

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
