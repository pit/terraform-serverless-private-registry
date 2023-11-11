module "lambda_providers_download" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-providers-download"
  description   = "Registry API: :namespace/:type/:version/download/:os/:arch"
  handler       = "registry-lambda-aws"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = {
    BUCKET_NAME            = aws_s3_bucket.this.id
    LAMBDA_TYPE            = "providers-download"
    PGP_KEY_ID             = "PGP_KEY_ID"
    PGP_KEY_PUBLIC_ARMORED = "PGP_KEY_PUBLIC_ARMORED"
  }

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  create_package = false
  s3_existing_package = {
    bucket = var.s3_bucket
    key    = var.s3_bucket_key
  }

  tags = merge({
    Name = "${var.name_prefix}-providers-download"
  }, var.tags)
}
