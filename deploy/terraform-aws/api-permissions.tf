resource "aws_lambda_permission" "lambdas_permissions" {
  for_each = local.api_routes

  function_name = each.value.lambda

  statement_id = "AllowInvokeFromApiGateway"
  action       = "lambda:InvokeFunction"
  principal    = "apigateway.amazonaws.com"

  source_arn = "${module.api.apigatewayv2_api_execution_arn}/*"
}
