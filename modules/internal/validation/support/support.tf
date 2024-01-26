variable "min_support" {
  type = string
  validation {
    condition     = var.min_support == "experimental" || var.min_support == "beta" || var.min_support == "ga"
    error_message = "min_support must be experimental, beta or ga"
  }
}

variable "use_case" {
  type = string
}

variable "use_case_support" {
  type = map(any)
}

locals {
  support = try(var.use_case_support[var.use_case], "unsupported")
}

check "support" {

  // min_support experimental can be satisfied by experimental, beta or ga
  assert {
    condition     = var.min_support != "experimental" || local.support == "experimental" || local.support == "beta" || local.support == "ga"
    error_message = "Support is ${local.support} but min_support is ${var.min_support}"
  }

  // min_support beta can be satisfied by beta or ga
  assert {
    condition     = var.min_support != "beta" || local.support == "beta" || local.support == "ga"
    error_message = "Support is ${local.support} but min_support is ${var.min_support}"
  }

  // min_support ga can be satisfied by only ga
  assert {
    condition     = var.min_support != "ga" || local.support == "ga"
    error_message = "Support is ${local.support} but min_support is ${var.min_support}"
  }
}
