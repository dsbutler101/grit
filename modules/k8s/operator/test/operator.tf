module "operator" {
  source = "../internal/"

  operator_version   = var.operator_version
  override_manifests = var.override_manifests
}
