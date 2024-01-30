variable "region" {
  description = "The GCP region in which your cluster will reside"
}

variable "zone" {
  description = "The GCP zone in which to create your cluster"
}

variable "gke_num_nodes" {
  description = "The number of GKE nodes in your cluster"
}
