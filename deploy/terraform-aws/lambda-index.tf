module "lambda_index" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-index"
  description   = "Registry API: /"
  handler       = "index"
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
    key    = aws_s3_bucket_object.lambda_index.key
  }

  tags = merge({
    Name = "${var.name_prefix}-index"
  }, var.tags)
}

data "archive_file" "lambda_index" {
  type        = "zip"
  source_file = "${var.distrib_dir}/index"
  output_path = "${path.module}/index.zip"
}

resource "aws_s3_bucket_object" "lambda_index" {
  bucket = var.s3_bucket
  key    = "${var.s3_prefix}terraform-serverless-private-registry/${var.lambda_version}/index.zip"
  source = data.archive_file.lambda_index.output_path
  etag   = filemd5("${var.distrib_dir}/index")
}
