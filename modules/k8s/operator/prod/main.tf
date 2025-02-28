#######################
# METADATA VALIDATION #
#######################

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "k8s-operator"
  use_case_support = tomap({
    "k8s-operator" = "experimental"
  })
  min_support = var.metadata.min_support
}

##########################
# KUBERNETES PROD MODULE #
##########################

module "operator" {
  source = "../internal/"

  operator_version   = var.operator_version
  override_manifests = var.override_manifests
}
