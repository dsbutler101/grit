locals {
  # Note: Do not write files under /tmp during boot 
  # https://cloudinit.readthedocs.io/en/latest/reference/modules.html#write-files
  cloudwatch_agent_json_path = "/var/tmp/amazon-cloudwatch-agent.json"

  cloud_config = {
    write_files = setunion([
      {
        path        = local.cloudwatch_agent_json_path
        owner       = "root:root"
        permissions = "0644"
        content     = base64decode(var.cloudwatch_agent_json)
      },
      ],
    var.node_exporter.enabled ? module.node_exporter[0].write_files_config : [])
    runcmd = local.runcmd
  }

  install_cloudwatch_agent_commands = [
    "sudo curl -O https://amazoncloudwatch-agent.s3.amazonaws.com/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb",
    "sudo apt-get -o DPkg::Lock::Timeout=300 install ./amazon-cloudwatch-agent.deb -y",
    "sudo usermod -aG adm cwagent",
    "sudo amazon-cloudwatch-agent-ctl -a start",
    "sudo amazon-cloudwatch-agent-ctl -a fetch-config -c file:${local.cloudwatch_agent_json_path} -s",
  ]

  runcmd = concat(
    local.install_cloudwatch_agent_commands,
    var.node_exporter.enabled ? module.node_exporter[0].commands : []
  )

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

module "node_exporter" {
  count                 = var.node_exporter.enabled ? 1 : 0
  source                = "../../../node_exporter"
  node_exporter_port    = var.node_exporter.port
  node_exporter_version = var.node_exporter.version
}
