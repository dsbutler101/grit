variable "project_id" {
  description = "The GCP project in which to create your cluster"
}

variable "region" {
  default     = "us-central1"
  description = "The GCP region in which your cluster will reside"
}

variable "zone" {
  default     = "us-central1-a"
  description = "The GCP zone in which to create your cluster"
}

variable "gke_num_nodes" {
  description = "The number of GKE nodes in your cluster"
  default     = 2
}
