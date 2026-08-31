resource "aws_sqs_queue" "diplomacy_api_events_dlq" {
  name = "diplomacy-api-v1-events-dlq"

  tags = {
    Application = "diplomacy-api"
  }
}

resource "aws_sqs_queue" "diplomacy_api_events" {
  name = "diplomacy-api-v1-events"

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.diplomacy_api_events_dlq.arn
    maxReceiveCount     = 5
  })

  tags = {
    Application = "diplomacy-api"
  }
}

resource "aws_iam_policy" "lambda_sqs" {
  name = "diplomacy-lambda-sqs-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "sqs:ChangeMessageVisibility",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:ReceiveMessage",
        "sqs:SendMessage"
      ]
      Resource = aws_sqs_queue.diplomacy_api_events.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_sqs_access" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.lambda_sqs.arn
}

output "events_queue_url" {
  value = aws_sqs_queue.diplomacy_api_events.url
}
