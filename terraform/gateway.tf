variable "aws_region" {
  type    = string
  default = "us-west-2"
}

locals {
  openapi = templatefile("openapi.json", {
    region     = var.aws_region
    turns_arn  = aws_lambda_function.fn["turns"].arn
    games_arn  = aws_lambda_function.fn["games"].arn
    maps_arn  = aws_lambda_function.fn["maps"].arn
    phases_arn  = aws_lambda_function.fn["phases"].arn
    players_arn  = aws_lambda_function.fn["players"].arn
    board_arn  = aws_lambda_function.fn["board"].arn
    orders_arn = aws_lambda_function.fn["orders"].arn
    test_arn   = aws_lambda_function.fn["test"].arn
  })
}

resource "aws_apigatewayv2_api" "http_api" {
  name          = "diplomacy-api"
  protocol_type = "HTTP"
  version       = "1.0.0"

  body = local.openapi
}

resource "aws_apigatewayv2_stage" "v1" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "v1"
  auto_deploy = true
}

resource "aws_lambda_permission" "allow_apigw_invoke" {
  for_each      = aws_lambda_function.fn
  statement_id  = "AllowExecutionFromHttpApi_${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = each.value.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

output "base_url" {
  value = "${aws_apigatewayv2_api.http_api.api_endpoint}/${aws_apigatewayv2_stage.v1.name}"
}
