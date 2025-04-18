output "grit_runner_managers" {
  value = tomap({
    "runner-manager" = tomap({
      instance_name   = module.runner.ec2_instance_name
      instance_id     = module.runner.ec2_instance_id
      address         = "${module.runner.ec2_public_ip}:22"
      wrapper_address = module.runner.ec2_runner_wrapper_socket_path
      username        = "ubuntu"
      ssh_key_pem     = module.runner.ec2_ssh_key_openssh_pem
    })
  })
  sensitive = true
}
