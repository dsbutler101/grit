locals {
  support      = try(var.use_case_support[var.use_case], "unsupported")
  fail_message = "Support for the '${var.use_case}' use case is ${local.support} but min_support is ${var.min_support}"
}

// min_support none can be satisfied by unsupported, experimental, beta or ga
module "check-min-support-none" {
  source  = "../fail_validation"
  message = var.min_support != "none" || local.support == "unsupported" || local.support == "experimental" || local.support == "beta" || local.support == "ga" ? "" : local.fail_message
}

// min_support experimental can be satisfied by experimental, beta or ga
module "check-min-support-experimental" {
  source  = "../fail_validation"
  message = var.min_support != "experimental" || local.support == "experimental" || local.support == "beta" || local.support == "ga" ? "" : local.fail_message
}

// min_support beta can be satisfied by beta or ga
module "check-min-support-beta" {
  source  = "../fail_validation"
  message = var.min_support != "beta" || local.support == "beta" || local.support == "ga" ? "" : local.fail_message
}

// min_support ga can be satisfied by only ga
module "check-min-support-ga" {
  source  = "../fail_validation"
  message = var.min_support != "ga" || local.support == "ga" ? "" : local.fail_message
}
