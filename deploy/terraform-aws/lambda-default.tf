module "lambda_default" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-default"
  description   = "Registry API: $default"
  handler       = "default"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  create_package = false
  s3_existing_package = {
    bucket = var.s3_bucket
    key    = aws_s3_bucket_object.lambda_default.key
  }

  tags = merge({
    Name = "${var.name_prefix}-default"
  }, var.tags)
}

data "archive_file" "lambda_default" {
  type        = "zip"
  source_file = "${var.distrib_dir}/default"
  output_path = "${path.module}/default.zip"
}

resource "aws_s3_bucket_object" "lambda_default" {
  bucket = var.s3_bucket
  key    = "${var.s3_prefix}terraform-serverless-private-registry/${var.lambda_version}/default.zip"
  source = data.archive_file.lambda_default.output_path
  etag   = filemd5("${var.distrib_dir}/default")
}
