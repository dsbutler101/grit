module "dev" {
  source = "../../modules/fleeting/aws"

  asg_storage = {
    type       = "gp3"
    throughput = 750
  }

  autoscaling_groups = {
    main = {
      ami_id = var.ami
      instance_type = "mac2.metal"
      subnet_cidr = "10.0.0.0/24"
    }
  }

  aws_vpc_cidr = "10.0.0.0/24"

  labels = var.labels
}