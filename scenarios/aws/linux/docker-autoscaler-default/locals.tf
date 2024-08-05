locals {
  default_labels = {
    managed = "grit"
  }

  metadata = {
    name        = var.name
    labels      = merge(var.labels, local.default_labels)
    min_support = "experimental"
  }

  ami_id = module.ami_lookup.ami_id
}
