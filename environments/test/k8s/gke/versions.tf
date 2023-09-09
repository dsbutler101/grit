# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.74.0"
    }

    gitlab = {
      source = "gitlabhq/gitlab"
      version = "16.3.0"
    }
  }

  required_version = ">= 0.14"
}

