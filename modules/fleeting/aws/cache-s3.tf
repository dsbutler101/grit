resource "aws_s3_bucket" "runners-cache" {
  bucket = var.cache_bucket_name

  tags = local.tags
}

resource "aws_s3_bucket_ownership_controls" "runners-cache" {
  bucket = aws_s3_bucket.runners-cache.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "runners-cache-acl" {
  depends_on = [
    aws_s3_bucket_ownership_controls.runners-cache
  ]

  bucket = aws_s3_bucket.runners-cache.id
  acl    = "private"
}

resource "aws_s3_bucket_lifecycle_configuration" "runners-cache-lifecycle" {
  bucket = aws_s3_bucket.runners-cache.id

  rule {
    id = "remove_files_older_than_14_days"

    status = "Enabled"

    expiration {
      days = 14
    }
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "runners-cache-sse" {
  bucket = aws_s3_bucket.runners-cache.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}