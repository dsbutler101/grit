config {
  call_module_type = "all"
}

# TODO: pick a sensible minimum tf version to apply
rule "terraform_required_version" {
  enabled = false
}

# TODO: enable snake_case naming checks
# will require "moving" many existing resources
rule "terraform_naming_convention" {
  enabled = false
}

rule "terraform_standard_module_structure" {
  enabled = true
}

rule "terraform_required_providers" {
  enabled = false
}
