package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

type (
	RequestDto struct {
		KeyId            string `json:"keyId"`
		Sha256Sums       string `json:"sha256Sums"`
		Sha256SumsSig    string `json:"sha256SumsSig"`
		Sha256SumsSigPub string `json:"sha256SumsSigPub"`
	}

	ResponseOk struct {
		Status string `json:"status"`
	}
)

func init() {
	logger, _ = helpers.InitLogger("DEBUG", true)
	bucketName := os.Getenv("BUCKET_NAME")
	storageSvc, _ := storage.NewStorage(bucketName, logger)
	providersSvc, _ = providers.NewProviders(storageSvc, logger)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer logger.Sync()
	logger.Debug(fmt.Sprintf("%s lambda called", request.RequestContext.RequestID),
		zap.Reflect("request", request),
	)

	providerNamespace := request.PathParameters["namespace"]
	providerType := request.PathParameters["type"]
	providerVersion := request.PathParameters["version"]

	//keyId := request.QueryStringParameters["keyId"]
	//sha256Sums := request.QueryStringParameters["sha256sums"]
	//sha256SumsSig := request.QueryStringParameters["sha256sumsSig"]
	//sha256SumsSigPub := request.QueryStringParameters["sha256sumsSigPub"]

	var dto RequestDto
	err := json.Unmarshal([]byte(request.Body), &dto)
	if err != nil {
		return helpers.ApiErrorUnknown(request.RequestContext.RequestID), nil
	}

	params := providers.SaveSignaturesInput{
		Namespace:        &providerNamespace,
		Type:             &providerType,
		Version:          &providerVersion,
		KeyId:            &dto.KeyId,
		Sha256Sums:       &dto.Sha256Sums,
		Sha256SumsSig:    &dto.Sha256SumsSig,
		Sha256SumsSigPub: &dto.Sha256SumsSigPub,
	}
	resp, providersErr := providersSvc.SaveSignatures(request.RequestContext.RequestID, params)

	if providersErr != nil {
		logger.Error(fmt.Sprintf("%s error %s", request.RequestContext.RequestID, err.Error()),
			zap.Error(err),
		)
		return helpers.ApiErrorUnknown(request.RequestContext.RequestID), nil
	}

	if resp.ProviderExists != nil && *resp.ProviderExists == false {
		msg := fmt.Sprintf("Provider %s/%s with version %s not found", providerNamespace, providerType, providerVersion)
		logger.Error(msg)
		return helpers.ApiErrorNotFound(request.RequestContext.RequestID, msg), nil
	}

	if resp.MetadataExists != nil && *resp.MetadataExists == true {
		msg := fmt.Sprintf("Provider %s/%s with version %s already has metadata", providerNamespace, providerType, providerVersion)
		logger.Error(msg)
		return helpers.ApiErrorConflict(request.RequestContext.RequestID, msg), nil
	}

	if resp.WrongContent != nil && *resp.WrongContent == true {
		msg := fmt.Sprintf("Wrong content found: %s", *resp.Details)
		logger.Error(msg)
		return helpers.ApiErrorBadRequest(request.RequestContext.RequestID, msg), nil
	}

	return helpers.ApiResponse(http.StatusOK, ResponseOk{Status: "ok"}), nil
}
