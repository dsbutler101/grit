module "test" {
  source           = "../../modules/test"
  manager_provider = "helm"
  runner_provider  = "gke"
}