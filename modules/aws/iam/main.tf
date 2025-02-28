#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}

###################
# IAM PROD MODULE #
###################

resource "aws_iam_user" "fleeting-service-account" {
  name = var.metadata.name

  tags = merge(var.metadata.labels, {
    Name = var.metadata.name
  })
}

# both `name` and `tags` are unsupported arguments
data "aws_iam_policy_document" "fleeting-service-account-policy-document" {
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
  name        = var.metadata.name
  description = "A policy for accessing autoscaling groups"
  policy      = data.aws_iam_policy_document.fleeting-service-account-policy-document.json

  tags = merge(var.metadata.labels, {
    Name = var.metadata.name
  })
}

# both `name` and `tags` are unsupported arguments
resource "aws_iam_user_policy_attachment" "fleeting-service-account-attach" {
  user       = aws_iam_user.fleeting-service-account.name
  policy_arn = aws_iam_policy.fleeting-service-account-policy.arn
}

# both `name` and `tags` are unsupported arguments
resource "aws_iam_access_key" "fleeting-service-account-key" {
  user = aws_iam_user.fleeting-service-account.name
}

