module "lambda_authorizer" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-authorizer"
  description   = "Registry API Authorizer"
  handler       = "authorizer"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = { for user_name, user_obj in var.users : "USER_${upper(md5(user_name))}" => md5(user_obj.password) }

  create_role = false
  lambda_role = module.lambdas_role.iam_role_arn

  attach_cloudwatch_logs_policy     = true
  cloudwatch_logs_retention_in_days = 7

  create_package = false
  s3_existing_package = {
    bucket = var.s3_bucket
    key    = aws_s3_bucket_object.lambda_authorizer.key

  }

  tags = merge({
    Name = "${var.name_prefix}-authorizer"
  }, var.tags)
}

data "archive_file" "lambda_authorizer" {
  type        = "zip"
  source_file = "${var.distrib_dir}/authorizer"
  output_path = "${path.module}/authorizer.zip"
}

resource "aws_s3_bucket_object" "lambda_authorizer" {
  bucket = var.s3_bucket
  key    = "${var.s3_prefix}terraform-serverless-private-registry/${var.lambda_version}/authorizer.zip"
  source = data.archive_file.lambda_authorizer.output_path
  etag   = filemd5("${var.distrib_dir}/authorizer")
}
