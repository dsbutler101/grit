resource "aws_iam_user" "cache" {
  name = "${var.metadata.name}-cache"

  tags = merge(var.metadata.labels, {
    Name = var.metadata.name
  })
}

data "aws_iam_policy_document" "cache_bucket_access_policy_document" {
  statement {
    effect = "Allow"

    resources = [
      "${aws_s3_bucket.cache.arn}/*"
    ]

    actions = [
      "s3:GetBucketLockation",
      "s3:ListBucket",
      "s3:PutObject",
      "s3:GetObject",
      "s3:DeleteObject",
      "s3:GetObjectVersion"
    ]
  }
}

resource "aws_iam_policy" "cache_bucket_access_policy" {
  name        = "${var.metadata.name}-cache"
  description = "A policy for accessing S3 cache bucket"
  policy      = data.aws_iam_policy_document.cache_bucket_access_policy_document.json

  tags = merge(var.metadata.labels, {
    Name = var.metadata.name
  })
}

resource "aws_iam_user_policy_attachment" "cache_bucket_user_policy" {
  user       = aws_iam_user.cache.name
  policy_arn = aws_iam_policy.cache_bucket_access_policy.arn
}

resource "aws_iam_access_key" "cache_bucket_user_key" {
  user = aws_iam_user.cache.name
}
