module "lambda_custom_providers_upload" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-custom-providers-upload"
  description   = "Registry API: GET /:namespace/:type/:version/:os/:arch/upload"
  handler       = "custom-providers-upload"
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
    key    = aws_s3_bucket_object.lambda_custom_providers_upload.key
  }

  tags = merge({
    Name = "${var.name_prefix}-custom-providers-upload"
  }, var.tags)
}

data "archive_file" "lambda_custom_providers_upload" {
  type        = "zip"
  source_file = "${var.distrib_dir}/custom-providers-upload"
  output_path = "${path.module}/custom-providers-upload.zip"
}

resource "aws_s3_bucket_object" "lambda_custom_providers_upload" {
  bucket = var.s3_bucket
  key    = "${var.s3_prefix}terraform-serverless-private-registry/${var.lambda_version}/custom-providers-upload.zip"
  source = data.archive_file.lambda_custom_providers_upload.output_path
  etag   = filemd5("${var.distrib_dir}/custom-providers-upload")
}
