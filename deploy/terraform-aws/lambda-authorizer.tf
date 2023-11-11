module "lambda_authorizer" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "v2.36.0"

  function_name = "${var.name_prefix}-authorizer"
  description   = "Registry API Authorizer"
  handler       = "registry-lambda-aws"
  runtime       = "go1.x"

  memory_size = 256
  timeout     = 5

  environment_variables = merge(
    {
      LAMBDA_TYPE = "authorizer"
    },
    { for user_name, user_obj in var.users : "USER_${upper(md5(user_name))}" => md5(user_obj.password) },
  )

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
    Name = "${var.name_prefix}-authorizer"
  }, var.tags)
}
