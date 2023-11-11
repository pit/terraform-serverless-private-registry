module "lambda_providers_download" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-providers-download"
  description   = "Registry API: :namespace/:type/:version/download/:os/:arch"
  handler       = "providers-download"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = {
    BUCKET_NAME            = aws_s3_bucket.this.id
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
    key    = aws_s3_bucket_object.lambda_providers_download.key
  }

  tags = merge({
    Name = "${var.name_prefix}-providers-download"
  }, var.tags)
}

data "archive_file" "lambda_providers_download" {
  type        = "zip"
  source_file = "${var.distrib_dir}/providers-download"
  output_path = "${path.module}/providers-download.zip"
}

resource "aws_s3_bucket_object" "lambda_providers_download" {
  bucket = var.s3_bucket
  key    = "${var.s3_prefix}terraform-serverless-private-registry/${var.lambda_version}/providers-download.zip"
  source = data.archive_file.lambda_providers_download.output_path
  etag   = filemd5("${var.distrib_dir}/providers-download")
}
