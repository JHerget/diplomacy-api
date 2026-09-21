data "aws_caller_identity" "current" {}

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

resource "aws_iam_policy" "scheduler_invoke_turns_lambda" {
  name = "diplomacy-turn-scheduler-lambda-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = "lambda:InvokeFunction"
      Resource = aws_lambda_function.fn["turns"].arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "scheduler_invoke_turns_lambda" {
  role       = aws_iam_role.scheduler_turns_role.name
  policy_arn = aws_iam_policy.scheduler_invoke_turns_lambda.arn
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
        Resource = "arn:aws:scheduler:${var.aws_region}:${data.aws_caller_identity.current.account_id}:schedule/default/diplomacy-game-*"
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

output "turn_scheduler_role_arn" {
  value = aws_iam_role.scheduler_turns_role.arn
}

output "turn_schedule_lambda_arn" {
  value = local.turns_lambda_arn
}
