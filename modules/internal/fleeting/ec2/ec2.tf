###############################
# INTERNAL EC2 INSTANCE GROUP #
###############################

module "common" {
  source = "./common"

  license_arn                      = try(module.macos.license-config-arn, "")
  jobs-host-resource-group-outputs = try(module.macos.jobs-host-resource-group-outputs, {})

  scale_min       = var.scale_min
  scale_max       = var.scale_max
  idle_percentage = var.idle_percentage

  asg_storage = {
    type       = "gp3"
    throughput = 750
  }
  autoscaling_groups = {
    main = {
      ami_id        = var.vm_img_id
      instance_type = var.instance_type
      subnet_cidr   = "10.0.0.0/24"
    }
  }
  aws_vpc_cidr = "10.0.0.0/24"
}

module "macos" {
  count  = var.os == "macos" ? 1 : 0
  source = "./macos"

  autoscaling_groups = {
    main = {
      ami_id        = var.vm_img_id
      instance_type = var.instance_type
      subnet_cidr   = "10.0.0.0/24"
    }
  }

  aws_vpc_cidr = "10.0.0.0/24"
}
