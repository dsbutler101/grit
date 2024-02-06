variable "name" {
  type = string
}

variable "google_region" {
  type = string
}

variable "subnetworks" {
  type = map(string)
}
