module "test" {
  source            = "../../modules/aws/test"
  manager_provider  = "ec2"
  fleeting_service  = "ec2"
  gitlab_project_id = var.gitlab_project_id
}