locals {
  bucket_name = var.bucket_name != "" ? var.bucket_name : "${var.metadata.name}-runner-cache"
}

resource "aws_s3_bucket" "cache_bucket_server_logs" {
  bucket = "${local.bucket_name}-logs"

  force_destroy = true

  tags = merge(var.metadata.labels, {
    Name = local.bucket_name
  })
}

resource "aws_s3_bucket_versioning" "cache_bucket_server_logs" {
  bucket = aws_s3_bucket.cache_bucket_server_logs.id

  versioning_configuration {
    status = "Disabled"
  }
}

resource "aws_s3_bucket_public_access_block" "cache_bucket_server_logs" {
  bucket = aws_s3_bucket.cache_bucket_server_logs.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket" "cache" {
  bucket = local.bucket_name

  force_destroy = true

  tags = merge(var.metadata.labels, {
    Name = local.bucket_name
  })
}

resource "aws_s3_bucket_versioning" "cache" {
  bucket = aws_s3_bucket.cache.id

  versioning_configuration {
    status = "Disabled"
  }
}

resource "aws_s3_bucket_public_access_block" "cache" {
  bucket = aws_s3_bucket.cache.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "cache" {
  bucket = aws_s3_bucket.cache.id

  rule {
    id     = "default"
    status = "Enabled"

    // Empty filter - apply to all objects in the bucket
    filter {}

    expiration {
      days = var.cache_object_lifetime
    }
  }
}

resource "aws_s3_bucket_logging" "cache" {
  bucket = aws_s3_bucket.cache.id

  target_bucket = aws_s3_bucket.cache_bucket_server_logs.id
  target_prefix = "logs/"

  target_object_key_format {
    partitioned_prefix {
      partition_date_source = "EventTime"
    }
  }
}
