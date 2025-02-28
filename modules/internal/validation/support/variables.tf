variable "min_support" {
  type = string
  validation {
    condition     = var.min_support == "none" || var.min_support == "experimental" || var.min_support == "beta" || var.min_support == "ga"
    error_message = "min_support must be none, experimental, beta or ga. A min_support of 'none' means that no support is required and should only be used in testing and development."
  }
}

variable "use_case" {
  type = string
}

variable "use_case_support" {
  type = map(any)
}
