resource "aws_launch_template" "fleeting_asg_template" {
  # Name must be alphanumeric, no spaces
  name = var.name

  description = "Launch template for GitLab Runner fleeting configuration"

  image_id      = var.ephemeral_runner_ami
  instance_type = var.instance_type

  key_name = aws_key_pair.jobs_key_pair.key_name

  user_data = var.install_cloudwatch_agent ? data.cloudinit_config.fleeting_config.rendered : null

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }

  dynamic "license_specification" {
    for_each = var.license_arn != "" ? [1] : []

    content {
      license_configuration_arn = var.license_arn
    }
  }

  dynamic "placement" {
    for_each = length(var.jobs_host_resource_group_outputs) > 0 ? [1] : []

    content {
      tenancy                 = "host"
      host_resource_group_arn = var.jobs-host-resource-group-outputs["ResourceGroupArn"]
    }
  }

  dynamic "iam_instance_profile" {
    for_each = [var.instance_role_profile_name]

    content {
      name = iam_instance_profile.value
    }
  }

  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      delete_on_termination = "true"
      volume_size           = var.storage_size
      volume_type           = var.storage_type
      throughput            = var.storage_throughput
      encrypted             = var.ebs_encryption
      kms_key_id            = var.kms_key_arn
    }
  }

  network_interfaces {
    delete_on_termination = true
    security_groups       = var.security_group_ids
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
resource "aws_autoscaling_group" "fleeting_asg" {
  name = var.name

  dynamic "launch_template" {
    for_each = var.mixed_instances_policy != null ? [] : [1]

    content {
      id      = aws_launch_template.fleeting_asg_template.id
      version = aws_launch_template.fleeting_asg_template.latest_version
    }
  }

  min_size = var.scale_min
  max_size = var.scale_max

  health_check_grace_period = 600

  vpc_zone_identifier = var.subnet_ids

  protect_from_scale_in = true

  dynamic "mixed_instances_policy" {
    for_each = var.mixed_instances_policy != null ? [var.mixed_instances_policy] : []

    content {
      launch_template {
        launch_template_specification {
          launch_template_id = aws_launch_template.fleeting_asg_template.id
          version            = aws_launch_template.fleeting_asg_template.latest_version
        }

        dynamic "override" {
          for_each = try(mixed_instances_policy.value.override, [])
          content {
            instance_type = try(override.value.instance_type, null)
          }
        }
      }
    }
  }

  dynamic "tag" {
    for_each = var.labels
    content {
      key                 = tag.key
      value               = tag.value
      propagate_at_launch = false
    }
  }
}
