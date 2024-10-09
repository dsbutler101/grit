variable "name" {
  type = string
}

variable "labels" {
  type = map(any)
}

variable "cache_object_lifetime" {
  type = number
}

variable "bucket_name" {
  type = string
}
