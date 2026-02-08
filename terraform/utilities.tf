resource "aws_secretsmanager_secret" "diplomacy_secrets" {
    name = "diplomacy-credentials"
    description = "Credentials for the Diplomacy API"

    recovery_window_in_days = 7

    tags = {
        Application = "diplomacy-api"
    }
}

resource "aws_s3_bucket" "maps_bucket" {
    bucket = "diplomacy-maps-01"

    tags = {
        Application = "diplomacy-api"
    }
}

resource "aws_s3_bucket_public_access_block" "maps_bucket" {
    bucket = aws_s3_bucket.maps_bucket.id

    block_public_acls = true
    block_public_policy = true
    ignore_public_acls = true
    restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "maps_bucket" {
    bucket = aws_s3_bucket.maps_bucket.id

    rule {
        apply_server_side_encryption_by_default {
            sse_algorithm = "AES256"
        }
    }
}
