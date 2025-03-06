variable "operator_version" {
  type        = string
  description = <<-EOF
    The operator version to deploy. This should be specified in semantic version format
    (e.g. 'v1.2.3') or set to 'latest' to use the most recent release.
  EOF
}

variable "override_manifests" {
  type        = string
  description = <<-EOT
    Optional path to custom operator manifests. Supports the following formats:
      - HTTP(S) URL (e.g., "https://example.com/custom-operator.yaml")
      - Local file path with "file://" prefix (e.g., "file:///path/to/operator.yaml")
      - If empty, uses the official GitLab Runner Operator manifest
  EOT

  validation {
    condition     = var.override_manifests == "" || can(regex("^(https?://|file://)", var.override_manifests))
    error_message = "override_manifests must be empty or start with 'http://', 'https://', or 'file://'"
  }
}
