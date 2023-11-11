gobuildcmd = GOSUMDB="off" GOPRIVATE=github.com/pit/terraform-serverless-private-registry CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  go build -installsuffix cgo -ldflags "-X main.version=`cat version` -X main.builddate=`date -u +.%Y%m%d.%H%M%S` -w -s"
version_snapshot = $(shell date '+%Y%m%d-%H%M%S')
# version_snapshot = local

.PHONY: dep
dep:
	go mod download

.PHONY: build-lambda
build_aws:
	mkdir -p bin/aws
	$(gobuildcmd) -o bin/aws/authorizer lambdas/aws/authorizer/*.go

# 	# lambda for index response
	$(gobuildcmd) -o bin/aws/default   lambdas/aws/default/*.go
	$(gobuildcmd) -o bin/aws/index     lambdas/aws/index/*.go
	$(gobuildcmd) -o bin/aws/discovery lambdas/aws/discovery/*.go

# 	# https://www.terraform.io/docs/internals/module-registry-protocol.html
# 	$(gobuildcmd) -o bin/aws/modules-list           lambdas/aws/modules-list/*.go
# 	$(gobuildcmd) -o bin/aws/modules-search         lambdas/aws/modules-search/*.go
	$(gobuildcmd) -o bin/aws/modules-versions       lambdas/aws/modules-versions/*.go
	$(gobuildcmd) -o bin/aws/modules-download       lambdas/aws/modules-download/*.go
# 	$(gobuildcmd) -o bin/aws/modules-latest-version lambdas/aws/modules-latest-version/*.go
# 	$(gobuildcmd) -o bin/aws/modules-get            lambdas/aws/modules-get/*.go
	$(gobuildcmd) -o bin/aws/custom-modules-upload lambdas/aws/custom-modules-upload/*.go

# 	# https://www.terraform.io/docs/internals/provider-registry-protocol.html
	$(gobuildcmd) -o bin/aws/providers-versions lambdas/aws/providers-versions/*.go
	$(gobuildcmd) -o bin/aws/providers-download lambdas/aws/providers-download/*.go
	$(gobuildcmd) -o bin/aws/custom-providers-upload lambdas/aws/custom-providers-upload/*.go
	$(gobuildcmd) -o bin/aws/custom-providers-checksums-upload lambdas/aws/custom-providers-checksums-upload/*.go

pack_aws:
	mkdir -p dist/aws
	echo $(version_snapshot) > dist/aws/version_snapshot
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-authorizer_$(version_snapshot).zip bin/aws/authorizer

	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-default_$(version_snapshot).zip   bin/aws/default
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-index_$(version_snapshot).zip     bin/aws/index
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-discovery_$(version_snapshot).zip bin/aws/discovery

# 	zip -j -u dist/modules-list.zip bin/modules-list
# 	zip -j -u dist/modules-search.zip bin/modules-search
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-modules-versions_$(version_snapshot).zip    bin/aws/modules-versions
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-modules-download_$(version_snapshot).zip    bin/aws/modules-download
# 	zip -j -u dist/modules-latest-version.zip bin/modules-latest-version
# 	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-modules-get_$(version_snapshot).zip         bin/aws/modules-get
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-custom-modules-upload_$(version_snapshot).zip  bin/aws/custom-modules-upload

	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-providers-versions_$(version_snapshot).zip  bin/aws/providers-versions
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-providers-download_$(version_snapshot).zip  bin/aws/providers-download
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-custom-providers-upload_$(version_snapshot).zip  bin/aws/custom-providers-upload
	zip -j -u dist/aws/terraform-serverless-private-registry_aws-lambda-custom-providers-checksums-upload_$(version_snapshot).zip  bin/aws/custom-providers-checksums-upload

test_aws:
	export BUCKET_NAME=terraform-registry-kvinta-io-dev
	export BASE_DIR=$(pwd)
	go test -i