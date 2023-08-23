#######################################################
# Maintainer access for Runner SaaS engineers         #
# assigned through our eng-dev-verify-runner-5442c67c #
# sandbox account                                     #
#######################################################

data "aws_iam_policy" "ec2_full_access" {
  name = "AmazonEC2FullAccess"
}

data "aws_iam_policy" "ecr_full_access" {
  name = "AmazonEC2ContainerRegistryFullAccess"
}

data "aws_iam_policy" "iam_full_access" {
  name = "IAMFullAccess"
}

data "aws_iam_policy" "s3_full_access" {
  name = "AmazonS3FullAccess"
}

data "aws_iam_policy" "vpc_full_access" {
  name = "AmazonVPCFullAccess"
}

data "aws_iam_policy" "service_quotas_full_access" {
  name = "ServiceQuotasFullAccess"
}

data "aws_iam_policy" "support_full_access" {
  name = "AWSSupportAccess"
}

data "aws_iam_policy_document" "resource_groups_full_access" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "resource-groups:*"
    ]
  }
}

resource "aws_iam_policy" "resource_groups_full_access" {
  name        = "resourceGroupsFullAccess"
  description = "A policy with full access to resource groups management"

  policy = data.aws_iam_policy_document.resource_groups_full_access.json
}

data "aws_iam_policy_document" "license_manager_full_access" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "license-manager:*",
    ]
  }
}

resource "aws_iam_policy" "license_manager_full_access" {
  name        = "licenseManagerFullAccess"
  description = "A policy with full access to AWS License Manager"

  policy = data.aws_iam_policy_document.license_manager_full_access.json
}

data "aws_iam_policy_document" "eng_dev_verify_runner" {
  statement {
    actions = [
      "sts:AssumeRole",
      "sts:TagSession",
      "sts:SetSourceIdentity",
    ]

    principals {
      type = "AWS"

      identifiers = [
        "arn:aws:iam::${local.eng_dev_verify_runner_account_id}:root"
      ]
    }
  }
}

resource "aws_iam_role" "eng_dev_verify_runner" {
  name               = "eng_dev_verify_runner"
  assume_role_policy = data.aws_iam_policy_document.eng_dev_verify_runner.json

  managed_policy_arns = [
    data.aws_iam_policy.ec2_full_access.arn,
    data.aws_iam_policy.ecr_full_access.arn,
    data.aws_iam_policy.iam_full_access.arn,
    data.aws_iam_policy.s3_full_access.arn,
    data.aws_iam_policy.vpc_full_access.arn,
    data.aws_iam_policy.service_quotas_full_access.arn,
    data.aws_iam_policy.support_full_access.arn,
    aws_iam_policy.resource_groups_full_access.arn,
    aws_iam_policy.license_manager_full_access.arn,
  ]
}
