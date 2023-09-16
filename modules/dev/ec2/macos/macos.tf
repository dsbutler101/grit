################################
# DEV MACOS EC2 INSTANCE GROUP #
################################

module "instance_group" {
  source = "../../../internal/fleeting/ec2/macos"
  asg_storage = {
    type       = "gp3"
    throughput = 750
  }
  autoscaling_groups = {
    main = {
      ami_id        = "ami-034ccb74da463ebe1"
      instance_type = "mac2.metal"
      subnet_cidr   = "10.0.0.0/24"
    }
  }
  aws_vpc_cidr = "10.0.0.0/24"
}