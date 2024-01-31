variable "name" {
  type = string
}

variable "labels" {
  type = map(string)
}

variable "cache_object_lifetime" {
  type = number
}

variable "bucket_name" {
  type = string
}

variable "bucket_location" {
  type = string
}

variable "service_account_emails" {
  type = list(string)
}
