package aws_custome_modules_upload

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"net/http"
	"os"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/modules"
	"terraform-serverless-private-registry/lib/storage"

	"github.com/Masterminds/semver/v3"
)

var (
	modulesSvc *modules.Modules
	logger     *zap.Logger
)

func Setup() {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ = helpers.InitLogger("DEBUG", true)
	logger.Debug("Lambda loading")

	storageSvc, _ := storage.NewStorage(bucketName, logger)
	modulesSvc, _ = modules.NewModules(storageSvc, logger)
}

func init() {
	Setup()
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

	_, err := semver.StrictNewVersion(version)
	if err != nil {
		logger.Error(fmt.Sprintf("%s error %s", request.RequestContext.RequestID, err.Error()),
			zap.Error(err),
		)
		return helpers.ApiErrorIncorrectVersion(request.RequestContext.RequestID, fmt.Sprintf("%s: %s", version, err.Error())), nil
	}

	params := modules.ModuleParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
		Version:   &version,
	}

	resp, err := modulesSvc.GetUploadUrl(request.RequestContext.RequestID, params)

	if err != nil {
		logger.Error(fmt.Sprintf("%s error %s", request.RequestContext.RequestID, err.Error()),
			zap.Error(err),
		)
		return helpers.ApiErrorUnknown(request.RequestContext.RequestID), nil
	}

	if resp.ModuleExists {
		msg := fmt.Sprintf("Module %s/%s/%s version %s already exists", namespace, name, provider, version)
		return helpers.ApiErrorConflict(request.RequestContext.RequestID, msg), nil
	}

	return helpers.ApiResponse(http.StatusOK, resp), nil
}
