resource "aws_iam_service_linked_role" "license-manager" {
  aws_service_name = "license-manager.amazonaws.com"
}

resource "aws_licensemanager_license_configuration" "license-config" {
  depends_on = [
    aws_iam_service_linked_role.license-manager
  ]

  name                     = "${var.shard}-required-license"
  license_count            = var.required_license_count_per_asg * var.cores_per_license * length(var.autoscaling_groups)
  license_count_hard_limit = false
  license_counting_type    = "Core"

  tags = local.tags
}

locals {
  jobs_host_resource_group = "${var.shard}-host-resource-group"
}

resource "aws_cloudformation_stack" "jobs-host-resource-group" {
  name = local.jobs_host_resource_group

  tags = local.tags

  template_body = <<EOS
    {
        "Resources" : {
            "HostResourceGroup": {
                "Type": "AWS::ResourceGroups::Group",
                "Properties": {
                    "Name": "${local.jobs_host_resource_group}",
                    "Configuration": [
                        {
                            "Type": "AWS::EC2::HostManagement",
                            "Parameters": [
                                {
                                    "Name": "allowed-host-based-license-configurations",
                                    "Values": [
                                        "${aws_licensemanager_license_configuration.license-config.arn}"
                                    ]
                                },
                                {
                                    "Name": "allowed-host-families",
                                    "Values": [
                                        "mac2"
                                    ]
                                },
                                {
                                    "Name": "auto-allocate-host",
                                    "Values": [
                                        "true"
                                    ]
                                },
                                {
                                    "Name": "auto-release-host",
                                    "Values": [
                                        "true"
                                    ]
                                }
                            ]
                        },
                        {
                            "Type": "AWS::ResourceGroups::Generic",
                            "Parameters": [
                                {
                                    "Name": "allowed-resource-types",
                                    "Values": [
                                        "AWS::EC2::Host"
                                    ]
                                },
                                {
                                    "Name": "deletion-protection",
                                    "Values": [
                                        "UNLESS_EMPTY"
                                    ]
                                }
                            ]
                        }
                    ]
                }
            }
        },
        "Outputs" : {
            "ResourceGroupArn" : {
                "Description": "ResourceGroup Arn",
                "Value" : {
                    "Fn::GetAtt": [
                        "HostResourceGroup",
                        "Arn"
                    ]
                },
                "Export" : {
                    "Name": "${local.jobs_host_resource_group}"
                }
            }
        }
    }
EOS
}