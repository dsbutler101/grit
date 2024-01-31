variable "runner_token" {
  default = "MY_RUNNER_TOKEN"
}

module "runner" {
  source = "../../modules/aws/runner/test"
  metadata = {
    name = "my-little-runner"
  }
  service = "ec2"
  gitlab = {
    runner_token = var.runner_token
  }
  vpc = {
    id        = "vpc-0d119da238d878eef"
    subnet_id = "subnet-0bd3ab8c221e14bfc"
  }
}
