config {
  call_module_type = "all"
}

# TODO: pick a sensible minimum tf version to apply
rule "terraform_required_version" {
  enabled = false
}

rule "terraform_naming_convention" {
  enabled = true
}

rule "terraform_standard_module_structure" {
  enabled = true
}

rule "terraform_required_providers" {
  enabled = false
}
