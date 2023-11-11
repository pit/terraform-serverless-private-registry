package aws_providers_versions

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

type (
	Response struct {
		Versions []providers.ProviderVersion `json:"versions"`
	}
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

	params := providers.ListProviderVersionsInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
	}
	resp, err := providersSvc.ListProviderVersions(request.RequestContext.RequestID, params)

	if err != nil {
		logger.Error(fmt.Sprintf("%s error %s", request.RequestContext.RequestID, err.Error()),
			zap.Error(err),
		)
		return helpers.ApiErrorUnknown(request.RequestContext.RequestID), nil
	}

	if resp.ProviderExists != nil && *resp.ProviderExists == false {
		return helpers.ApiErrorNotFound(request.RequestContext.RequestID, fmt.Sprintf("Provider %s/%s doesn't exists", providerNamespace, providerType)), nil
	}

	return helpers.ApiResponse(http.StatusOK, Response{Versions: resp.Versions}), nil
}
