module "api_domain_acm" {
  source  = "terraform-aws-modules/acm/aws"
  version = "v3.3.0"

  count = var.create_domain ? 1 : 0

  domain_name = var.record_name == "" ? var.domain_name : "${var.record_name}.${var.domain_name}"
  zone_id     = data.aws_route53_zone.domain.zone_id

  validation_method    = "DNS"
  validate_certificate = true
  wait_for_validation  = true

  validation_allow_overwrite_records = true

  tags = merge(var.tags, {
    Name = "terraform-registry-v2-dev"
  })
}

data "aws_route53_zone" "domain" {
  name = var.domain_name
}

module "api_domain_route53" {
  source  = "terraform.kvinta.io/infra/route53/aws"
  version = "1.6.0"

  count = var.create_domain ? 1 : 0

  domain  = var.domain_name
  zone_id = data.aws_route53_zone.domain.zone_id

  aliases = [{
    name    = var.record_name
    zone_id = aws_apigatewayv2_domain_name.base[0].domain_name_configuration[0].hosted_zone_id
    value   = aws_apigatewayv2_domain_name.base[0].domain_name_configuration[0].target_domain_name
  }]
}
