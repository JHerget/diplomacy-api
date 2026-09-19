data "aws_caller_identity" "current" {}

locals {
  default_event_bus_arn = "arn:aws:events:${var.aws_region}:${data.aws_caller_identity.current.account_id}:event-bus/default"
}

data "aws_iam_policy_document" "scheduler_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["scheduler.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "scheduler_turns_role" {
  name               = "diplomacy-turn-scheduler-role"
  assume_role_policy = data.aws_iam_policy_document.scheduler_assume_role.json
}

resource "aws_iam_policy" "scheduler_put_turn_events" {
  name = "diplomacy-turn-scheduler-events-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = "events:PutEvents"
      Resource = local.default_event_bus_arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "scheduler_put_turn_events" {
  role       = aws_iam_role.scheduler_turns_role.name
  policy_arn = aws_iam_policy.scheduler_put_turn_events.arn
}

data "aws_iam_policy_document" "eventbridge_api_target_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["events.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "eventbridge_api_target_role" {
  name               = "diplomacy-eventbridge-api-target-role"
  assume_role_policy = data.aws_iam_policy_document.eventbridge_api_target_assume_role.json
}

resource "aws_iam_policy" "eventbridge_invoke_turns_api" {
  name = "diplomacy-eventbridge-turns-api-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = "execute-api:Invoke"
      Resource = "${aws_apigatewayv2_api.http_api.execution_arn}/${aws_apigatewayv2_stage.v1.name}/POST/games/*/turns"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "eventbridge_invoke_turns_api" {
  role       = aws_iam_role.eventbridge_api_target_role.name
  policy_arn = aws_iam_policy.eventbridge_invoke_turns_api.arn
}

resource "aws_cloudwatch_event_rule" "scheduled_turns" {
  name = "diplomacy-api-scheduled-turns"

  event_pattern = jsonencode({
    source        = ["diplomacy-api.turn-scheduler"]
    "detail-type" = ["CreateTurn"]
  })
}

resource "aws_cloudwatch_event_target" "scheduled_turns_api" {
  rule     = aws_cloudwatch_event_rule.scheduled_turns.name
  arn      = "${aws_apigatewayv2_api.http_api.execution_arn}/${aws_apigatewayv2_stage.v1.name}/POST/games/*/turns"
  role_arn = aws_iam_role.eventbridge_api_target_role.arn
  input    = "{}"

  http_target {
    path_parameter_values = ["$.detail.gameId"]
  }
}

resource "aws_iam_policy" "lambda_turn_scheduler" {
  name = "diplomacy-lambda-turn-scheduler-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "scheduler:CreateSchedule",
          "scheduler:DeleteSchedule",
          "scheduler:GetSchedule"
        ]
        Resource = "arn:aws:scheduler:${var.aws_region}:${data.aws_caller_identity.current.account_id}:schedule/default/diplomacy-turn-*"
      },
      {
        Effect   = "Allow"
        Action   = "iam:PassRole"
        Resource = aws_iam_role.scheduler_turns_role.arn
        Condition = {
          StringEquals = {
            "iam:PassedToService" = "scheduler.amazonaws.com"
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_turn_scheduler_access" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.lambda_turn_scheduler.arn
}
