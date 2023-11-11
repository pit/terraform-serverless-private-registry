resource "aws_apigatewayv2_stage" "base" {
  api_id      = module.api.apigatewayv2_api_id
  name        = "base"
  auto_deploy = true

  stage_variables = {
    LOG_LEVEL = "DEBUG"
  }

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.access_logs.arn
    format          = var.access_log_format
  }

  default_route_settings {
    detailed_metrics_enabled = true
    throttling_burst_limit   = 5000
    throttling_rate_limit    = 10000
  }

  dynamic "route_settings" {
    for_each = { for key in keys(local.api_routes) : key => local.api_routes[key] if key != "$default" }
    content {
      route_key                = route_settings.key
      detailed_metrics_enabled = true
      throttling_burst_limit   = 5000
      throttling_rate_limit    = 10000
    }
  }

  tags = merge(var.tags, {
    Name = "terraform-registry/base"
  })

  lifecycle {
    ignore_changes = [deployment_id]
  }
}

resource "aws_apigatewayv2_domain_name" "base" {
  count = var.create_domain ? 1 : 0

  domain_name = var.record_name == "" ? var.domain_name : "${var.record_name}.${var.domain_name}"

  domain_name_configuration {
    certificate_arn = module.api_domain_acm[0].acm_certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = merge({
    Name = "terraform-registry/base"
  }, var.tags)
}

resource "aws_apigatewayv2_api_mapping" "base" {
  count = var.create_domain ? 1 : 0

  api_id      = module.api.apigatewayv2_api_id
  domain_name = aws_apigatewayv2_domain_name.base[0].id
  stage       = aws_apigatewayv2_stage.base.id
}

resource "aws_cloudwatch_log_group" "access_logs" {
  name = "/terraform-registry/${var.record_name == "" ? var.domain_name : "${var.record_name}.${var.domain_name}"}"

  retention_in_days = 7

  tags = merge(var.tags, {
    Name = "terraform-registry/base"
  })
}
