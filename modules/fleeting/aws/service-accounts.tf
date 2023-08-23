#########################################
# AWS Autoscaling Group service account #
#########################################

resource "aws_iam_user" "fleeting-service-account" {
  name = "fleeting-service-account"
}

data "aws_iam_policy_document" "fleeting-service-account-policy" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "autoscaling:DescribeAutoScalingInstances",
      "autoscaling:DescribeAutoScalingGroups",
      "ec2:DescribeInstances"
    ]
  }

  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "autoscaling:SetDesiredCapacity",
      "autoscaling:TerminateInstanceInAutoScalingGroup"
    ]
  }

  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "ec2:GetPasswordData",
      "ec2-instance-connect:SendSSHPublicKey"
    ]
  }
}

resource "aws_iam_policy" "fleeting-service-account-policy" {
  name        = "fleeting-service-account-policy"
  description = "A policy for accessing autoscaling groups"
  policy      = data.aws_iam_policy_document.fleeting-service-account-policy.json
}

resource "aws_iam_user_policy_attachment" "fleeting-service-account-attach" {
  user       = aws_iam_user.fleeting-service-account.name
  policy_arn = aws_iam_policy.fleeting-service-account-policy.arn
}

resource "aws_iam_access_key" "fleeting-service-account-key" {
  user = aws_iam_user.fleeting-service-account.name
}

#############################
# AWS cache service account #
#############################

resource "aws_iam_user" "s3-cache-service-account" {
  name = "s3-cache-service-account"
}

data "aws_iam_policy_document" "s3-cache-service-account-policy" {
  statement {
    effect = "Allow"
    resources = [
      "${aws_s3_bucket.runners-cache.arn}/*",
      aws_s3_bucket.runners-cache.arn
    ]

    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:ListBucket"
    ]
  }
}

resource "aws_iam_policy" "s3-cache-service-account-policy" {
  name        = "s3-cache-service-account-policy"
  description = "A policy for accessing S3 cache"
  policy      = data.aws_iam_policy_document.s3-cache-service-account-policy.json
}

resource "aws_iam_user_policy_attachment" "s3-cache-attach" {
  user       = aws_iam_user.s3-cache-service-account.name
  policy_arn = aws_iam_policy.s3-cache-service-account-policy.arn
}

resource "aws_iam_access_key" "cache-service-account-key" {
  user = aws_iam_user.s3-cache-service-account.name
}
