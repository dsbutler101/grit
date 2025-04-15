#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}

##################
# DEFAULT LABELS #
##################

module "labels" {
  source = "../../internal/labels"

  name              = var.metadata.name
  additional_labels = var.metadata.labels
}

###################
# IAM PROD MODULE #
###################

resource "aws_iam_user" "fleeting_service_account" {
  name = var.metadata.name

  tags = module.labels.merged
}

# both `name` and `tags` are unsupported arguments
data "aws_iam_policy_document" "fleeting_service_account_policy_document" {
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

resource "aws_iam_policy" "fleeting_service_account_policy" {
  name        = var.metadata.name
  description = "A policy for accessing autoscaling groups"
  policy      = data.aws_iam_policy_document.fleeting_service_account_policy_document.json

  tags = module.labels.merged
}

# both `name` and `tags` are unsupported arguments
resource "aws_iam_user_policy_attachment" "fleeting_service_account_attach" {
  user       = aws_iam_user.fleeting_service_account.name
  policy_arn = aws_iam_policy.fleeting_service_account_policy.arn
}

# both `name` and `tags` are unsupported arguments
resource "aws_iam_access_key" "fleeting_service_account_key" {
  user = aws_iam_user.fleeting_service_account.name
}

