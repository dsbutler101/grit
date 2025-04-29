terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    google = {
      source = "hashicorp/google"
    }
    gitlab = {
      source = "gitlabhq/gitlab"
    }
    kubectl = {
      source = "alekc/kubectl"
    }
    http = {
      source = "hashicorp/http"
    }
    tls = {
      source = "hashicorp/tls"
    }
    cloudinit = {
      source = "hashicorp/cloudinit"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}
