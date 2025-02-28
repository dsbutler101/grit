resource "aws_licensemanager_license_configuration" "license-config" {
  name                     = var.name
  license_count            = var.license_count_per_asg * var.cores_per_license
  license_count_hard_limit = false
  license_counting_type    = "Core"

  tags = merge(var.labels, {
    Name = var.name
  })
}

locals {
  jobs_host_resource_group = var.name
}

resource "aws_cloudformation_stack" "jobs-cloudformation-stack" {
  name = local.jobs_host_resource_group

  tags = merge(var.labels, {
    Name = var.name
  })

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

