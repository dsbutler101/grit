locals {
  default_labels = {
    managed = "grit"
  }

  name_label = {
    name = var.name
  }

  labels = merge(local.default_labels, local.name_label, var.additional_labels)
}

module "validate_labels" {
  source = "../../internal/validation/labels"
  labels = local.labels
}