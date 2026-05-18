locals {
    lambdas = {
        "board" = {
            handler = "lambdas.board"
        }
        "orders" = {
            handler = "lambdas.orders"
        }
        "test" = {
            handler = "lambdas.test"
        }
    }
}

data "archive_file" "lambda_zip" {
    type = "zip"
    source_dir = "../dist"
    output_path = "../dist/lambdas.zip"
}

data "aws_iam_policy_document" "lambda_assume_role" {
    statement {
        actions = ["sts:AssumeRole"]
        principals {
            type = "Service"
            identifiers = ["lambda.amazonaws.com"]
        }
    }
}

resource "aws_iam_policy" "lambda_read_secret" {
    name = "diplomacy-lambda-secrets-policy"

    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [{
            Effect   = "Allow"
            Action   = [
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret"
            ]
            Resource = aws_secretsmanager_secret.diplomacy_secrets.arn
        }]
    })
}

resource "aws_iam_policy" "lambda_s3" {
    name = "diplomacy-lambda-s3-policy"

    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = [
                  "s3:GetObject",
                  "s3:PutObject",
                  "s3:DeleteObject"
                ]
                Resource = "${aws_s3_bucket.maps_bucket.arn}/*"
            },
            {
                Effect = "Allow"
                Action = "s3:ListBucket"
                Resource = aws_s3_bucket.maps_bucket.arn
            }
        ]
    })
}

resource "aws_iam_role" "lambda_role" {
    name = "diplomacy-lambda-role"
    assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy_attachment" "lambda_basic_logs" {
    role = aws_iam_role.lambda_role.name
    policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# resource "aws_iam_role_policy_attachment" "lambda_vpc_access" {
#     role = aws_iam_role.lambda_role.name
#     policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
# }

resource "aws_iam_role_policy_attachment" "lambda_secret_access" {
    role = aws_iam_role.lambda_role.name
    policy_arn = aws_iam_policy.lambda_read_secret.arn
}

resource "aws_iam_role_policy_attachment" "lambda_s3_access" {
    role = aws_iam_role.lambda_role.name
    policy_arn = aws_iam_policy.lambda_s3.arn
}

# resource "aws_security_group" "lambda_vpc" {
#     name = "diplomacy-lambda-vpc"
#     vpc_id = aws_vpc.main.id
#
#     egress {
#         from_port = 0
#         to_port = 0
#         protocol = "-1"
#         cidr_blocks = ["0.0.0.0/0"]
#     }
# }

resource "aws_lambda_function" "fn" {
    for_each = local.lambdas

    filename = data.archive_file.lambda_zip.output_path
    source_code_hash = data.archive_file.lambda_zip.output_base64sha256

    function_name = "diplomacy-api-v1-${each.key}"
    role = aws_iam_role.lambda_role.arn
    handler = each.value.handler
    runtime = "nodejs24.x"
    memory_size = 1024
    timeout = 120

    # vpc_config {
    #     subnet_ids = [aws_subnet.private[0].id, aws_subnet.private[1].id]
    #     security_group_ids = [aws_security_group.lambda_vpc.id]
    # }
}
