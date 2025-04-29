terraform {
  required_version = ">= 0.13"

  required_providers {
    http = {
      source  = "hashicorp/http"
      version = "~> 3.4"
    }
  }
}

data "http" "gitlab_tags" {
  url = "https://gitlab.com/api/v4/projects/${var.project_id}/repository/tags?order_by=${var.order_by}"

  request_headers = {
    Accept = "application/json"
  }
}
