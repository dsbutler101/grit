## variables should be defined here and used in main.tf

variable "aws_zone" {
  type = string
}

variable "capacity_per_instance" {
  type = number

  default = 1
}

variable "max_instances" {
  type = number

  default = 200
}

variable "ephemeral_runner" {
  type = object({
    disk_type    = optional(string, "")
    disk_size    = optional(number, 25)
    machine_type = optional(string, "t2.medium")
    source_image = optional(string, "ami-0735db9b38fcbdb39")
  })

    default = {
    disk_type    = ""
    disk_size    = 25
    machine_type = "t2.medium"
    source_image = "ami-0735db9b38fcbdb39"
  }
}