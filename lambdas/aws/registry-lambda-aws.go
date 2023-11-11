package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	aws_authorizer "terraform-serverless-private-registry/lambdas/aws/authorizer"
	aws_custom_modules_upload "terraform-serverless-private-registry/lambdas/aws/custom-modules-upload"
	aws_custom_providers_checksums_upload "terraform-serverless-private-registry/lambdas/aws/custom-providers-checksums-upload"
	aws_custom_providers_upload "terraform-serverless-private-registry/lambdas/aws/custom-providers-upload"
	aws_default "terraform-serverless-private-registry/lambdas/aws/default"
	aws_discovery "terraform-serverless-private-registry/lambdas/aws/discovery"
	aws_index "terraform-serverless-private-registry/lambdas/aws/index"
	aws_modules_download "terraform-serverless-private-registry/lambdas/aws/modules-download"
	aws_modules_latest_version "terraform-serverless-private-registry/lambdas/aws/modules-latest-version"
	aws_modules_list "terraform-serverless-private-registry/lambdas/aws/modules-list"
	aws_modules_search "terraform-serverless-private-registry/lambdas/aws/modules-search"
	aws_modules_versions "terraform-serverless-private-registry/lambdas/aws/modules-versions"
	aws_providers_download "terraform-serverless-private-registry/lambdas/aws/providers-download"
	aws_providers_versions "terraform-serverless-private-registry/lambdas/aws/providers-versions"
	"terraform-serverless-private-registry/lib/helpers"
)

var handler interface{}

func init() {
	lambdaType := helpers.GetLambdaType()
	if lambdaType == helpers.LambdaTypeAuthorizer {
		handler = aws_authorizer.Handler
	} else if lambdaType == helpers.LambdaTypeDiscovery {
		handler = aws_discovery.Handler

	} else if lambdaType == helpers.LambdaTypeDefault {
		handler = aws_default.Handler
	} else if lambdaType == helpers.LambdaTypeIndex {
		handler = aws_index.Handler

	} else if lambdaType == helpers.LambdaTypeModulesDownload {
		handler = aws_modules_download.Handler
	} else if lambdaType == helpers.LambdaTypeModulesLatestVersion {
		handler = aws_modules_latest_version.Handler
	} else if lambdaType == helpers.LambdaTypeModulesList {
		handler = aws_modules_list.Handler
	} else if lambdaType == helpers.LambdaTypeModulesSearch {
		handler = aws_modules_search.Handler
	} else if lambdaType == helpers.LambdaTypeModulesVersions {
		handler = aws_modules_versions.Handler

	} else if lambdaType == helpers.LambdaTypeProvidersDownload {
		handler = aws_providers_download.Handler
	} else if lambdaType == helpers.LambdaTypeProvidersVersions {
		handler = aws_providers_versions.Handler

	} else if lambdaType == helpers.LambdaTypeCustomModulesUpload {
		handler = aws_custom_modules_upload.Handler
	} else if lambdaType == helpers.LambdaTypeCustomProvidersChecksumsUpload {
		handler = aws_custom_providers_checksums_upload.Handler
	} else if lambdaType == helpers.LambdaTypeCustomProvidersUpload {
		handler = aws_custom_providers_upload.Handler
	}
}

func main() {
	lambda.Start(handler)
}
