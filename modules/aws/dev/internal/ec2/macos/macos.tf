################################
# DEV MACOS EC2 INSTANCE GROUP #
################################

module "instance_group" {
  source = "../../../../internal/ec2/fleeting"
  os     = "macos"

  vm_img_id     = "ami-0fcd5ff1c92b00231"
  instance_type = "mac2.metal"
  aws_vpc_cidr  = "10.0.0.0/24"
}