locals {
  cloud_config = {
    write_files = [
      {
        path        = "/tmp/amazon-cloudwatch-agent.json"
        owner       = "root:root"
        permissions = "0644"
        content     = base64decode(var.cloudwatch_agent_json)
      }
    ]
    runcmd = local.cmd
  }

  cmd = [
    "sudo curl -O https://amazoncloudwatch-agent.s3.amazonaws.com/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb",
    "sudo apt-get -o DPkg::Lock::Timeout=300 install ./amazon-cloudwatch-agent.deb -y",
    "sudo usermod -aG adm cwagent",
    "sudo amazon-cloudwatch-agent-ctl -a start",
    "sudo amazon-cloudwatch-agent-ctl -a fetch-config -c file:/tmp/amazon-cloudwatch-agent.json -s"
  ]

  rendered_yaml = yamlencode(local.cloud_config)
}

data "cloudinit_config" "fleeting_config" {
  gzip          = false
  base64_encode = true

  part {
    filename     = "cloud-config.yaml"
    content_type = "text/cloud-config"

    content = local.rendered_yaml
  }
}
