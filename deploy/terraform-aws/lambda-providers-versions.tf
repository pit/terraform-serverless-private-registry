module "lambda_providers_versions" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-providers-versions"
  description   = "Registry API: :namespace/:type/versions"
  handler       = "providers-versions"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = {
    BUCKET_NAME = aws_s3_bucket.this.id
  }

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  create_package = false
  s3_existing_package = {
    bucket = var.s3_bucket
    key    = aws_s3_bucket_object.lambda_providers_versions.key
  }

  tags = merge({
    Name = "${var.name_prefix}-providers-versions"
  }, var.tags)
}

data "archive_file" "lambda_providers_versions" {
  type        = "zip"
  source_file = "${var.distrib_dir}/providers-versions"
  output_path = "${path.module}/providers-versions.zip"
}

resource "aws_s3_bucket_object" "lambda_providers_versions" {
  bucket = var.s3_bucket
  key    = "${var.s3_prefix}terraform-serverless-private-registry/${var.lambda_version}/providers-versions"
  source = data.archive_file.lambda_providers_versions.output_path
  etag   = filemd5("${var.distrib_dir}/providers-versions")
}
