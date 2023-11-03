module "dev" {
  source           = "../../modules/aws/dev"
  fleeting_service = "ec2"
  fleeting_os      = "macos"
}