resource "aws_launch_template" "fleeting-asg" {
  for_each = var.autoscaling_groups

  description = "Launch template for GitLab Runner fleeting configuration for ${each.key}"

  tags = local.tags

  image_id      = each.value.ami_id
  instance_type = each.value.instance_type

  key_name = aws_key_pair.jobs.key_name

  license_specification {
    license_configuration_arn = aws_licensemanager_license_configuration.license-config.arn
  }

  placement {
    tenancy                 = "host"
    host_resource_group_arn = aws_cloudformation_stack.jobs-host-resource-group.outputs["ResourceGroupArn"]
  }

  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      delete_on_termination = "true"
      volume_size           = var.asg_storage.size
      volume_type           = var.asg_storage.type
      throughput            = var.asg_storage.throughput
    }
  }

  network_interfaces {
    subnet_id = aws_subnet.jobs-vpc-subnet[each.key].id

    security_groups = [
      aws_security_group.jobs-security-group.id
    ]
  }

  tag_specifications {
    resource_type = "instance"

    tags = local.tags
  }

  lifecycle {
    ignore_changes = [
      tags,
      tag_specifications,
    ]
  }
}

resource "aws_autoscaling_group" "fleeting-asg" {
  for_each = var.autoscaling_groups

  name = each.key

  launch_template {
    id      = aws_launch_template.fleeting-asg[each.key].id
    version = aws_launch_template.fleeting-asg[each.key].latest_version
  }

  min_size = 0
  max_size = var.required_license_count_per_asg

  health_check_grace_period = 600

  vpc_zone_identifier = [
    aws_subnet.jobs-vpc-subnet[each.key].id
  ]

  protect_from_scale_in = var.protect_from_scale_in

  dynamic "tag" {
    for_each = local.tags
    content {
      key                 = tag.key
      value               = tag.value
      propagate_at_launch = false
    }
  }
}
