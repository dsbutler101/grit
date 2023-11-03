variable "gitlab_url" {}
variable "runner_token" {}
variable "aws_asg_name" {}
variable "idle_count" {}
variable "scale_max" {}
variable "executor" {}
variable "ssh_key_pem" {
  default = ""
}
variable "ssh_key_pem_name" {
  default = ""
}
variable "fleeting_service_account_secret_access_key" {
  default = ""
}
variable "fleeting_service_account_access_key_id" {
  default = ""
}
variable "fleeting_service" {
  description = "The system which provides infrastructure for the Runners"
}