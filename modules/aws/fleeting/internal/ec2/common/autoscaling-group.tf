resource "aws_launch_template" "fleeting-asg-template" {
  # Name must be alphanumeric, no spaces
  name = var.name

  description = "Launch template for GitLab Runner fleeting configuration"

  image_id      = var.ami_id
  instance_type = var.instance_type

  key_name = aws_key_pair.jobs-key-pair.key_name

  dynamic "license_specification" {
    for_each = var.license_arn != "" ? [1] : []

    content {
      license_configuration_arn = var.license_arn
    }
  }

  dynamic "placement" {
    for_each = length(var.jobs-host-resource-group-outputs) > 0 ? [1] : []

    content {
      tenancy                 = "host"
      host_resource_group_arn = var.jobs-host-resource-group-outputs["ResourceGroupArn"]
    }
  }

  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      delete_on_termination = "true"
      volume_size           = var.storage_size
      volume_type           = var.storage_type
      throughput            = var.storage_throughput
    }
  }

  network_interfaces {
    subnet_id = var.subnet_id

    security_groups = [
      aws_security_group.jobs-security-group.id
    ]
  }

  tag_specifications {
    resource_type = "instance"

    tags = merge(var.labels, {
      Name = var.name
    })
  }

  lifecycle {
    ignore_changes = [
      tags,
      tag_specifications,
    ]
  }

  tags = merge(var.labels, {
    Name = var.name
  })
}

# tags are a deprecated property on this resource type
resource "aws_autoscaling_group" "fleeting-asg" {
  name = var.name

  launch_template {
    id      = aws_launch_template.fleeting-asg-template.id
    version = aws_launch_template.fleeting-asg-template.latest_version
  }

  min_size = 0
  max_size = var.scale_max

  health_check_grace_period = 600

  vpc_zone_identifier = [
    var.subnet_id
  ]

  protect_from_scale_in = var.protect_from_scale_in

  dynamic "tag" {
    for_each = var.labels
    content {
      key                 = tag.key
      value               = tag.value
      propagate_at_launch = false
    }
  }
}
