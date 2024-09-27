resource "aws_iam_user" "cache" {
  name = "${var.name}-cache"

  tags = merge(var.labels, {
    Name = var.name
  })
}

data "aws_iam_policy_document" "cache-bucket-access-policy-document" {
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

resource "aws_iam_policy" "cache-bucket-access-policy" {
  name        = "${var.name}-cache"
  description = "A policy for accessing S3 cache bucket"
  policy      = data.aws_iam_policy_document.cache-bucket-access-policy-document.json

  tags = merge(var.labels, {
    Name = var.name
  })
}

resource "aws_iam_user_policy_attachment" "cache-bucket-user-policy" {
  user       = aws_iam_user.cache.name
  policy_arn = aws_iam_policy.cache-bucket-access-policy.arn
}

resource "aws_iam_access_key" "cache-bucket-user-key" {
  user = aws_iam_user.cache.name
}
