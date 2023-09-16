module "dev" {
  source            = "../../modules/dev"
  fleeting_provider = "ec2"
  os                = "macos"
}