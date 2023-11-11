locals {
  api_routes = {
    "GET /" = {
      lambda              = module.lambda_index.lambda_function_name
      authorizer_required = false
    },

    "GET /.well-known/terraform.json" = {
      lambda              = module.lambda_discovery.lambda_function_name
      authorizer_required = true
    },

    # TODO Phase 3
    # "GET /v1/modules" = {
    #   lambda = module.lambda_modules_list.lambda_function_name
    # },
    # TODO Phase 3
    # "GET /v1/modules/{namespace}" = {
    #   lambda = module.lambda_modules_list.lambda_function_name
    # },

    # TODO Phase 3
    # "GET /v1/modules/search" = {
    #   lambda = module.lambda_modules_search.lambda_function_name
    # },

    "GET /v1/modules/{namespace}/{name}/{provider}/versions" = {
      lambda              = module.lambda_modules_versions.lambda_function_name
      authorizer_required = true
    },

    # Unknown URL
    # "GET /v1/modules/{namespace}/{name}/{provider}/download" = {
    #   lambda = module.lambda_modules_download.lambda_function_name
    # },

    "GET /v1/modules/{namespace}/{name}/{provider}/{version}/download" = {
      lambda              = module.lambda_modules_download.lambda_function_name
      authorizer_required = true
    },
    # TODO Phase 2
    # TODO Implement module archive/metainfo upload
    # "POST /modules/custom/{namespace}/{name}/{provider}/{version}/upload" = {
    #   lambda = module.lambda_modules_upload.lambda_function_name
    # },

    # TODO Phase 3
    # "GET /v1/modules/{namespace}/{name}" = {
    #   lambda = module.lambda_modules_latest_version.lambda_function_name
    # },
    # TODO Phase 3
    # "GET /v1/modules/{namespace}/{name}/{provider}" = {
    #   lambda = module.lambda_modules_latest_version.lambda_function_name
    # },

    # TODO Phase 3
    # "GET /v1/modules/{namespace}/{name}/{provider}/{version}" = {
    #   lambda = module.lambda_modules_get.lambda_function_name
    # },


    "GET /v1/providers/{namespace}/{type}/versions" = {
      lambda              = module.lambda_providers_versions.lambda_function_name
      authorizer_required = true
    },

    "GET /v1/providers/{namespace}/{type}/{version}/download/{os}/{arch}" = {
      lambda              = module.lambda_providers_download.lambda_function_name
      authorizer_required = true
    },
    # TODO Phase 2
    # TODO Implement module/provider upload
    "GET /v1/custom/modules/{namespace}/{name}/{provider}/{version}/upload" = {
      lambda              = module.lambda_custom_modules_upload.lambda_function_name
      authorizer_required = true
    },
    "GET /v1/custom/providers/{namespace}/{type}/{version}/{os}/{arch}/upload" = {
      lambda              = module.lambda_custom_providers_upload.lambda_function_name
      authorizer_required = true
    },
    "POST /v1/custom/providers/{namespace}/{type}/{version}/checksums/upload" = {
      lambda              = module.lambda_custom_providers_checksums_upload.lambda_function_name
      authorizer_required = true
    },

    "$default" = {
      lambda = module.lambda_default.lambda_function_name
    },
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

module "api" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "v1.5.1"

  name          = var.name_prefix
  description   = "Serverless terraform registry API"
  protocol_type = "HTTP"

  # credentials_arn = module.apigateway_role.iam_role_arn
  # Custom domain
  # domain_name                      = var.domain_name
  # domain_name_certificate_arn      = var.domain_acm_arn
  create_api_domain_name           = false
  create_default_stage_api_mapping = true

  #TODO
  default_stage_access_log_destination_arn = "arn:aws:logs:eu-central-1:463422107539:log-group:/aws/api/artifacts/dev"
  default_stage_access_log_format          = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId"

  # Routes and integrations
  integrations = {
    for route_path in keys(local.api_routes) : (route_path) => merge(
      {
        lambda_arn             = "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:${lookup(local.api_routes[route_path], "lambda")}"
        credentials            = module.apigateway_role.iam_role_arn
        payload_format_version = "2.0"
        timeout_milliseconds   = 5000
        integration_type       = "AWS_PROXY"
      },
      can(local.api_routes[route_path].authorizer_required) ? local.api_routes[route_path].authorizer_required ? {
        authorization_type = "CUSTOM"
        authorizer_id      = aws_apigatewayv2_authorizer.basic_auth.id
      } : {} : {},
    )
  }

  tags = merge(var.tags, {
    Name = "terraform-registry"
  })
}
