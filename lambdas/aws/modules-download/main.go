package aws_modules_download

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
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
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ = helpers.InitLogger("DEBUG", true)
	logger.Debug("Lambda loading")

	storageSvc, _ := storage.NewStorage(bucketName, logger)
	modulesSvc, _ = modules.NewModules(storageSvc, logger)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()
	logger.Debug(fmt.Sprintf("%s lambda called", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	namespace := request.PathParameters["namespace"]
	name := request.PathParameters["name"]
	provider := request.PathParameters["provider"]
	version := request.PathParameters["version"]

	params := modules.ModuleParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
		Version:   &version,
	}
	resp, err := modulesSvc.GetDownloadUrl(request.RequestContext.RequestID, params)

	if err != nil {
		logger.Error(fmt.Sprintf("%s error %s", request.RequestContext.RequestID, err.Error()),
			zap.Error(err),
		)
		return helpers.ApiErrorUnknown(request.RequestContext.RequestID), nil
	}

	if resp.ModuleExists == false {
		return helpers.ApiErrorNotFound(request.RequestContext.RequestID, fmt.Sprintf("Module %s/%s/%s version %s doesn't exists", namespace, name, provider, version)), nil
	}

	result := helpers.ApiErrorNoContent()
	result.Headers["X-Terraform-Get"] = *resp.Url

	return result, nil
}
