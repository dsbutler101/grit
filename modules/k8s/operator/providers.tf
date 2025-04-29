terraform {
  required_version = ">= 0.13"

  required_providers {
    kubectl = {
      source  = "alekc/kubectl"
      version = "~> 2.0"
    }

    http = {
      source  = "hashicorp/http"
      version = "~> 3.4"
    }
  }
}
